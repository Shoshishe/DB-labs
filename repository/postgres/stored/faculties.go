package stored

import "github.com/google/uuid"

type Faculty struct {
	Id           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	UniversityId uuid.UUID `db:"university_id"`
}

