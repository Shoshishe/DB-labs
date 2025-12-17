package stored

import "github.com/google/uuid"

type University struct {
	Id        uuid.UUID `db:"uni_id"`
	Name      string    `db:"uni_name"`
	Shorthand string    `db:"uni_shorthand"`
}

func NewUniversity(id uuid.UUID, name, shorthand string) *University {
	return &University{Id: id, Name: name, Shorthand: shorthand}
}
