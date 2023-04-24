package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userid"
)

type contextKey string

const userIDKey contextKey = "userID"

func (h *Handler) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header, ok := r.Header["Authorization"]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID, err := h.authService.ParseToken(header[0])
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) RequireAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postId, err := strconv.Atoi(vars["id"])
		if err != nil || postId < 1 {
			newErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		userid := GetUserID(r)
		if userid < 0 {
			newErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		err = h.problemService.AccessRight(userid, postId)
		if err != nil {
			newErrorResponse(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) int {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0
	}
	return userID
}
