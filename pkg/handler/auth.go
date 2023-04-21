package handler

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"github.com/go-playground/validator"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var input models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validator := validator.New()
	err = validator.Struct(input)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUserService(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

type signinInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	var input signinInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validator := validator.New()
	err = validator.Struct(input)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateTokenService(input.Username, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jwtToken := models.Token{
		TokenString: token,
	}

	err = json.NewEncoder(w).Encode(jwtToken)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
