package responses

import "github.com/google/uuid"

type GroupsResponse struct {
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	ItemsCount  int8 `json:"count_per_page"`
	Groups      []struct {
		Id        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		FacultyId string    `json:"faculty_id"`
	} `json:"groups"`
}
