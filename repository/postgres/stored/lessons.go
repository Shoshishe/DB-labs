package stored

import "github.com/google/uuid"

type GroupGrades struct {
	UserId uuid.UUID `db:"user_id"`
	Grade  float64   `db:"group_grade"`
}
