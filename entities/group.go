package entities

import "github.com/google/uuid"

type Group struct {
	id           uuid.UUID
	name         string
	universityId uuid.UUID
}

func NewGroup(id uuid.UUID, name string, universityId uuid.UUID) *Group {
	return &Group{id: id, name: name, universityId: universityId}
}

func (gr *Group) Id() uuid.UUID {
	return gr.id
}

func (gr *Group) Name() string {
	return gr.name
}

func (gr *Group) UniversityId() uuid.UUID {
	return gr.universityId
}

