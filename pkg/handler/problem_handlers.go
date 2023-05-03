package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	userid := GetUserID(r)
	if userid == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var mergestruct models.ProblemWithTopics

	err := json.NewDecoder(r.Body).Decode(&mergestruct)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = mergestruct.Validate()
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.problemService.CreateProblem(userid, mergestruct)
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

func (h *Handler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.problemService.DeleteProblem(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(statusResponse{
		Status: "OK",
	})
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	var update models.ProblemUpdate
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.problemService.UpdateProblem(id, update)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(statusResponse{
		Status: "OK",
	})
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) ReadProblem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pageNum, err := strconv.Atoi(vars["pageNum"])
	if err != nil || pageNum < 1 {
		newErrorResponse(w, http.StatusBadRequest, "Invalid pageNum value")
		return
	}

	pageSize, err := strconv.Atoi(vars["pageSize"])
	if err != nil || pageNum < 1 {
		newErrorResponse(w, http.StatusBadRequest, "Invalid pageNum value")
		return
	}

	mergestruct, err := h.problemService.GetProblems(pageNum, pageSize)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(mergestruct)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) ReadProblemByParameter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	level := vars["level"]
	mergestruct, err := h.problemService.GetByParameter(topic, level)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(mergestruct)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) ReadProblemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	mergestruct, err := h.problemService.GetProblemByID(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(mergestruct)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) SearchProblem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	mergestruct, err := h.problemService.SearchProblem(title)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(mergestruct)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
