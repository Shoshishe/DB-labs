package entities

import "github.com/google/uuid"

type UserRole int8

const (
	Student UserRole = iota + 1
	Teacher
	Monitor
)

func RolesFromId(ids []int8) []UserRole {
	roles := []UserRole{}
	for _, id := range ids {
		if id >= int8(Student) && id <= int8(Monitor) {
			roles = append(roles, UserRole(id))
		}
	}
	return roles
}

func (role UserRole) String() string {
	switch role {
	case Student:
		return "Student"
	case Teacher:
		return "Teacher"
	case Monitor:
		return "Monitor"
	default:
		return ""
	}
}

func (usr *User) RoleNames() []string {
	res := []string{}
	for _, role := range usr.Roles() {
		res = append(res, role.String())
	}
	return res
}

type User struct {
	id           uuid.UUID
	universityId uuid.UUID
	roles        []UserRole
	name         string
	surname      string
	patronymic   string
	email        string
	password     string
}

func NewUser(id, universityId uuid.UUID, roles []UserRole, name, surname, patronymic, email, password string) *User {
	return &User{
		universityId: universityId,
		id:           id,
		name:         name,
		surname:      surname,
		patronymic:   patronymic,
		password:     password,
		email:        email,
		roles:        roles,
	}
}

func (usr *User) Id() uuid.UUID {
	return usr.id
}

func (usr *User) UniversityId() uuid.UUID {
	return usr.universityId
}

func (usr *User) Roles() []UserRole {
	return usr.roles
}

func (usr *User) Name() string {
	return usr.name
}

func (usr *User) Patronymic() string {
	return usr.patronymic
}

func (usr *User) Email() string {
	return usr.email
}

func (usr *User) Surname() string {
	return usr.surname
}

func (usr *User) Password() string {
	return usr.password
}
