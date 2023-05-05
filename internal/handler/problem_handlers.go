package handler

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
)

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	userid := GetUserID(r)
	if userid == 0 {

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	mergestruct := models.ProblemWithTopics{}
	mergestruct.Problem.UserID = userid

	if err := json.NewDecoder(r.Body).Decode(&mergestruct); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.problemService.CreateProblem(r.Context(), mergestruct)
	if err != nil {
		errorDefine(w, err)
		return
	}

	encodeJSONResponse(w, id, 200)
}

func (h *Handler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	params := models.URLParams{}
	if err := parseURLParams(r, &params); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.problemService.DeleteProblem(r.Context(), params.ID); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encodeJSONResponse(w, statusResponse{Status: "OK"}, 200)
}

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	params := models.URLParams{}
	if err := parseURLParams(r, &params); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	update := models.ProblemUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.problemService.UpdateProblem(r.Context(), params.ID, update); err != nil {
		errorDefine(w, err)
		return
	}

	encodeJSONResponse(w, statusResponse{Status: "OK"}, 200)
}

func (h *Handler) ReadProblem(w http.ResponseWriter, r *http.Request) {
	params := models.URLParams{}
	if err := parseURLParams(r, &params); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	mergestruct, err := h.problemService.GetProblems(r.Context(), params.PageNum, params.PageSize)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encodeJSONResponse(w, mergestruct, 200)
}

func (h *Handler) ReadProblemByParameter(w http.ResponseWriter, r *http.Request) {
	params := models.URLParams{}
	if err := parseURLParams(r, &params); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	mergestruct, err := h.problemService.GetByParameter(r.Context(), params.Topic, params.Level)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encodeJSONResponse(w, mergestruct, 200)
}

func (h *Handler) ReadProblemByID(w http.ResponseWriter, r *http.Request) {
	params := models.URLParams{}
	if err := parseURLParams(r, &params); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	mergestruct, err := h.problemService.GetProblemByID(r.Context(), params.ID)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encodeJSONResponse(w, mergestruct, 200)
}

func (h *Handler) SearchProblem(w http.ResponseWriter, r *http.Request) {
	params := models.URLParams{}
	if err := parseURLParams(r, &params); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	mergestruct, err := h.problemService.SearchProblem(r.Context(), params.Title)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encodeJSONResponse(w, mergestruct, 200)
}
