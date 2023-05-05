package handler

import (
	"git.01.alem.school/bbaktyke/test.project.git/cache"
	"git.01.alem.school/bbaktyke/test.project.git/internal/service"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type Handler struct {
	authService    service.AuthorizationService
	problemService service.ProblemService
	problemCache   cache.ProblemCash
	ch             *amqp.Channel
	qName          string
}

func NewHandler(authService service.AuthorizationService, problemService service.ProblemService, cache cache.ProblemCash, ch *amqp.Channel, qName string) *Handler {
	return &Handler{
		authService:    authService,
		problemService: problemService,
		problemCache:   cache,
		ch:             ch,
		qName:          qName,
	}
}

func (h *Handler) InitProblemRouters(router *mux.Router) {
	reqAuth := router.PathPrefix("").Subrouter()
	reqAccess := router.PathPrefix("").Subrouter()

	reqAuth.HandleFunc("/problem", h.CreateProblem).Methods("POST")
	reqAuth.HandleFunc("/problem/{pageNum}/{pageSize}", h.ReadProblem).Methods("GET")
	reqAuth.HandleFunc("/parameter/{topic}/{level}", h.ReadProblemByParameter).Methods("GET")
	reqAuth.HandleFunc("/search/{title}", h.SearchProblem).Methods("GET")
	reqAuth.HandleFunc("/problem/{id}", h.ReadProblemByID).Methods("GET")

	reqAccess.HandleFunc("/problem/{id}", h.UpdateProblem).Methods("PUT")
	reqAccess.HandleFunc("/problem/{id}", h.DeleteProblem).Methods("DELETE")

	reqAuth.Use(h.RequireAuthentication)
	reqAccess.Use(h.RequireAuthentication, h.RequireAccess)
}

func (h *Handler) InitAuthRouters(router *mux.Router) {
	router.HandleFunc("/signup", h.Signup).Methods("POST")
	router.HandleFunc("/signin", h.Signin).Methods("POST")
}
