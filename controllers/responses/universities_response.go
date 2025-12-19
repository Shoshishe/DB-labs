package responses

import (
	"db_labs/entities"

	"github.com/google/uuid"
)

type UniversityResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Shorthand string    `json:"shorthand"`
}

func UniversityFromEntity(uni *entities.University) *UniversityResponse {
	return &UniversityResponse{
		Id:        uni.Id(),
		Name:      uni.Name(),
		Shorthand: uni.Shorthand(),
	}
}
