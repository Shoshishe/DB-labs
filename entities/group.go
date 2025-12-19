package entities

import "github.com/google/uuid"

type Group struct {
	id           uuid.UUID
	name         string
	facultyId uuid.UUID
}

func NewGroup(id uuid.UUID, name string, facultyId uuid.UUID) *Group {
	return &Group{id: id, name: name, facultyId: facultyId}
}

func (gr *Group) Id() uuid.UUID {
	return gr.id
}

func (gr *Group) Name() string {
	return gr.name
}

func (gr *Group) FacultyId() uuid.UUID {
	return gr.facultyId
}

