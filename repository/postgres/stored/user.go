package stored

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	Id           uuid.UUID     `db:"id"`
	UniversityId uuid.UUID     `db:"university_id"`
	Roles        pq.Int64Array `db:"roles" json:",omitempty"`
	Name         string        `db:"name"`
	Surname      string        `db:"surname"`
	Patronymic   string        `db:"patronymic"`
	Email        string        `db:"email"`
	Password     string        `db:"password"`
}
