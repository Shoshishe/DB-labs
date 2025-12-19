package services

import (
	"context"
	"db_labs/controllers"
	"db_labs/entities"
	"fmt"
)

type UsersRepository interface {
	GetRoles(ctx context.Context) ([]int8, error)
	UpdateUser(ctx context.Context, request controllers.UpdateUserRequest) error
}

type UsersService struct {
	repo UsersRepository
}

func NewUsersService(repo UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (srv *UsersService) GetRoles(ctx context.Context) ([]entities.UserRole, error) {
	ids, err := srv.repo.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles in users service: %w", err)
	}
	roles := entities.RolesFromId(ids)
	return roles, nil
}

func (srv *UsersService) UpdateUser(ctx context.Context, request controllers.UpdateUserRequest) error {
	err := srv.repo.UpdateUser(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
