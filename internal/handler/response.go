package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.01.alem.school/bbaktyke/test.project.git/internal/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(errorResponse{
		Message: message,
	})
	w.WriteHeader(statusCode)
	w.Write(response)
}

func errorDefine(w http.ResponseWriter, err error) {
	var statusCode int
	switch err.Error() {
	case "update struct has no values":
		statusCode = 400

	case "title cannot be empty or contain only whitespaces":
		statusCode = 400
	case "description cannot be empty or contain only whitespaces":
		statusCode = 400
	case "level cannot be empty or contain only whitespaces":
		statusCode = 400
	case "samples cannot be empty or contain only whitespaces":
		statusCode = 400
	case "Not Valid values for level":
		statusCode = 400
	case "not valid topic id":
		statusCode = 400
	default:
		statusCode = 500
	}
	newErrorResponse(w, statusCode, err.Error())
}

func encodeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func parseURLParams(r *http.Request, params *models.URLParams) error {
	vars := mux.Vars(r)

	if id, err := strconv.Atoi(vars["id"]); err == nil {
		params.ID = id
	}

	params.Topic = vars["topic"]
	params.Level = vars["level"]

	if pageNum, err := strconv.Atoi(vars["pageNum"]); err == nil {
		params.PageNum = pageNum
	}

	if pageSize, err := strconv.Atoi(vars["pageSize"]); err == nil {
		params.PageSize = pageSize
	}

	params.Title = vars["title"]

	return nil
}
