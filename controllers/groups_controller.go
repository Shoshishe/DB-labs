package controllers

import (
	"context"
	"db_labs/controllers/constants"
	"db_labs/entities"
	"db_labs/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type GroupsService interface {
	GetGroups(ctx context.Context, itemsPerPage uint8, currentPage uint) ([]entities.Group, error)
}

type GroupsController struct {
	srv GroupsService
}

func NewGroupsController(srv GroupsService) *GroupsController {
	return &GroupsController{srv: srv}
}

func (controller *GroupsController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/groups", useTimeout(constants.DefaultTimeout, useAuthorized(http.HandlerFunc(controller.GetGroups))))
}

type GetGroupsRequest struct {
	ItemsPerPage uint8 `json:"per_page"`
	CurrentPage  uint  `json:"current_page"`
}

func (controller *GroupsController) GetGroups(w http.ResponseWriter, r *http.Request) {
	groupsRequest := &GetGroupsRequest{}
	err := json.NewDecoder(r.Body).Decode(groupsRequest)
	if err != nil {
		utils.JSONError(w, fmt.Sprintf("Failed to decode json body: %v", err), http.StatusBadRequest)
		return
	}
	groups, err := controller.srv.GetGroups(r.Context(), groupsRequest.ItemsPerPage, groupsRequest.CurrentPage)
	if err != nil {
		utils.JSONError(w, fmt.Sprintf("Failed to get groups: %v", err), http.StatusInternalServerError)
		slog.Error("Failed to get groups", "err", err)
		return
	}

	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		slog.Error("Failed to encode groups result", "err", err)
	}
}

type GetGroupGradesRequest struct {
	GroupId uuid.UUID `json:"group_id"`
}

func (controller *GroupsController) GetGradesByGroup(w http.ResponseWriter, r *http.Request) {
	groupGradesRequest := &GetGroupGradesRequest{}
	err := json.NewDecoder(r.Body).Decode(groupGradesRequest)
	if err != nil {
		utils.JSONError(w, fmt.Sprintf("Failed to decode json body: %v", err), http.StatusInternalServerError)
		return
	}

}
