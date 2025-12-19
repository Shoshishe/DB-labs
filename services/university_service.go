package services

import (
	"context"
	"db_labs/entities"
	"db_labs/repository/postgres/stored"
	"fmt"
)

type UniversityRepository interface {
	Select(ctx context.Context) ([]stored.University, error)
}

type UniversityService struct {
	repo UniversityRepository
}

func NewUniversityService(repo UniversityRepository) *UniversityService {
	return &UniversityService{repo: repo}
}

func (serv *UniversityService) GetUniversities(ctx context.Context) ([]entities.University, error) {
	stored, err := serv.repo.Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all universities in university service: %w", err)
	}
	universities := []entities.University{}
	for _, store := range stored {
		universities = append(universities, *entities.NewUniversity(store.Id, store.Name, store.Shorthand))
	}
	return universities, nil
}
