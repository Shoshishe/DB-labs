package entities

import "github.com/google/uuid"

type UserRole int8

const (
	Student UserRole = iota + 1
	Teacher
	Monitor
)

type User struct {
	id         uuid.UUID
	roles      []UserRole
	name       string
	surname    string
	patronymic string
	email      string
}
