package stored

import "github.com/google/uuid"

type Group struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	FacultyId uuid.UUID `db:"faculty_id"`
}

type SkippedHours struct {
	GroupName    string `db:"group_name"`
	SkippedHours int    `db:"skipped_hours"`
}
