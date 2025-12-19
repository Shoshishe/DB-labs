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

	"github.com/google/uuid"
)

type UserService interface {
	GetRoles(ctx context.Context) ([]entities.UserRole, error)
	UpdateUser(ctx context.Context, request UpdateUserRequest) error
}

type UserController struct {
	srv UserService
}

func NewUserController(srv UserService) *UserController {
	return &UserController{srv: srv}
}

func (controller *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /roles", useTimeout(constants.DefaultTimeout, http.HandlerFunc(controller.getRoles)))
	mux.Handle("PUT /user", useTimeout(constants.DefaultTimeout, useAuthorized(http.HandlerFunc(controller.updateUser))))
}

func (controller *UserController) getRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := controller.srv.GetRoles(r.Context())
	if err != nil {
		utils.JSONError(w, "Failed to get database roles", http.StatusInternalServerError)
		return
	}
	response := []responses.Role{}
	for _, role := range roles {
		response = append(response, *responses.NewRole(int8(role), role.String()))
	}
	json.NewEncoder(w).Encode(response)
}

type UpdateUserRequest struct {
	UniversityId uuid.UUID `json:"-"`
	Id           uuid.UUID `json:"-"`
	Name         *string   `json:"name,omitempty"`
	Surname      *string   `json:"surname,omitempty"`
	Patronymic   *string   `json:"patronymic,omitempty"`
	Roles        *[]int8   `json:"roles,omitempty"`
}

func (controller *UserController) updateUser(w http.ResponseWriter, r *http.Request) {
	var request UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.JSONError(w, fmt.Sprintf("Failed to unmarshal json request: %v", err), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(userIdCtxKey).(uuid.UUID)
	if !ok {
		utils.JSONError(w, "Failed to get credentials", http.StatusBadRequest)
		return
	}
	request.Id = userId
	universityId, ok := r.Context().Value(userIdCtxKey).(uuid.UUID)
	if !ok {
		utils.JSONError(w, "Failed to get credentials", http.StatusBadRequest)
		return
	}
	request.UniversityId = universityId

	err = controller.srv.UpdateUser(r.Context(), request)
	if err != nil {
		utils.JSONError(w, "Failed to update user info", http.StatusInternalServerError)
		slog.Error("Failed to update user info", "err", err.Error())
		return
	}
	err = json.NewEncoder(w).Encode("ok")
	if err != nil {
		slog.Error("Failed to encode ok response", "err", err.Error())
	}
}
