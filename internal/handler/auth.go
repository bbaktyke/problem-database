package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
	"github.com/go-playground/validator"
	"github.com/streadway/amqp"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.authService.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// publish message to RabbitMQ queue
	email, err := json.Marshal(input.Username)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        email,
	}
	if err := h.ch.Publish("", h.qName, false, false, msg); err != nil {
		log.Printf("failed to publish message to RabbitMQ queue: %v", err)
	}

	if err = json.NewEncoder(w).Encode(id); err != nil {
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
	if err = validator.Struct(input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jwtToken := models.Token{
		TokenString: token,
	}

	if err = json.NewEncoder(w).Encode(jwtToken); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
