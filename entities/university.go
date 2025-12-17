package entities

import "github.com/google/uuid"

type University struct {
	id        uuid.UUID
	name      string
	shorthand string
}

func NewUniversity(id uuid.UUID, name, shorthand string) *University {
	return &University{id: id, name: name, shorthand: shorthand}
}

func (uni *University) Id() uuid.UUID {
	return uni.id
}

func (uni *University) Name() string {
	return uni.name
}

func (uni *University) Shorthand() string {
	return uni.shorthand
}
