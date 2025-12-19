package services

import (
	"context"
	"db_labs/entities"
	"db_labs/repository/postgres/stored"
	"fmt"
)

type GroupsRepository interface {
	GetGroups(ctx context.Context, itemsPerPage uint8, currentPage uint) ([]stored.Group, error)
}

type GroupsService struct {
	repo GroupsRepository
}

func NewGroupsService(repo GroupsRepository) *GroupsService {
	return &GroupsService{repo: repo}
}

func (srv *GroupsService) GetGroups(ctx context.Context, itemsPerPage uint8, currentPage uint) ([]entities.Group, error) {
	stored, err := srv.repo.GetGroups(ctx, itemsPerPage, currentPage)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups from db: %w", err)
	}
	groups := []entities.Group{}
	for _, store := range stored {
		groups = append(groups, *entities.NewGroup(store.Id, store.Name, store.FacultyId))
	}
	return groups, nil
}
