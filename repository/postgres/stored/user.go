package stored

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID `db:"id"`
	UniversityId uuid.UUID `db:"university_id"`
	Roles        []int8	`db:"-"`
	Name         string `db:"name"`
	Surname      string `db:"surname"`
	Patronymic   string `db:"patronymic"`
	Email        string `db:"email"`
	Password     string `db:"password"`
}
