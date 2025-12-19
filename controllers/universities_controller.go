package controllers

import (
	"context"
	"db_labs/controllers/constants"
	"db_labs/controllers/responses"
	"db_labs/entities"
	"db_labs/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type UniversityController struct {
	serv UniversityService
}

type UniversityService interface {
	GetUniversities(ctx context.Context) ([]entities.University, error)
}

func NewUniversityController(mux *http.ServeMux, serv UniversityService) *UniversityController {
	return &UniversityController{serv: serv}
}

func (controller *UniversityController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /universities", useTimeout(constants.DefaultTimeout, (http.HandlerFunc(controller.getAll))))
}

func (controller *UniversityController) getAll(wr http.ResponseWriter, r *http.Request) {
	universities, err := controller.serv.GetUniversities(r.Context())
	if err != nil {
		utils.JSONError(wr, fmt.Sprintf("Failed to get universities: %v", err), http.StatusInternalServerError)
		return
	}
	response := []responses.UniversityResponse{}
	for _, uni := range universities {
		response = append(response, *responses.UniversityFromEntity(&uni))
	}
	err = json.NewEncoder(wr).Encode(response)
	if err != nil {
		slog.Error("Failed to encode universities get all response: ")
	}
}
