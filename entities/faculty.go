package entities

import "github.com/google/uuid"

type Faculty struct {
	id   uuid.UUID
	name string
}

func NewFaculty(id uuid.UUID, name string) *Faculty {
	return &Faculty{id: id, name: name}
}

func (fc *Faculty) Id() uuid.UUID {
	return fc.id
}

func (fc *Faculty) Name() string {
	return fc.name
}
