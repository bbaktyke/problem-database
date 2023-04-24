package handler

import (
	"git.01.alem.school/bbaktyke/test.project.git/pkg/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	authService    service.AuthorizationService
	problemService service.ProblemService
}

func NewHandler(authService service.AuthorizationService, problemService service.ProblemService) *Handler {
	return &Handler{authService: authService, problemService: problemService}
}

func (h *Handler) InitRouters() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/signup", h.Signup).Methods("POST")
	router.HandleFunc("/signin", h.Signin).Methods("POST")

	s := router.PathPrefix("").Subrouter()
	x := router.PathPrefix("").Subrouter()
	s.HandleFunc("/problem", h.CreateProblem).Methods("POST")
	s.HandleFunc("/problem/{pageNum}/{pageSize}", h.ReadProblem).Methods("GET")
	s.HandleFunc("/parameter/{topic}/{level}", h.ReadProblemByParameter).Methods("GET")
	s.HandleFunc("/search/{title}", h.SearchProblem).Methods("GET")

	s.HandleFunc("/problem/{id}", h.ReadProblemByID).Methods("GET")
	x.HandleFunc("/problem/{id}", h.UpdateProblem).Methods("PUT")
	x.HandleFunc("/problem/{id}", h.DeleteProblem).Methods("DELETE")
	s.Use(h.RequireAuthentication)
	x.Use(h.RequireAuthentication, h.RequireAccess)
	return router
}
