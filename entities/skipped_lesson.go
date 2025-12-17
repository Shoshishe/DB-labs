package entities

import "github.com/google/uuid"

type Skipped struct {
	studentId    uuid.UUID
	groupId      uuid.UUID
	skippedHours uint64
	lessonId     uuid.UUID
	isLegit      bool
}

func NewSkipped(studentId, groupId, lessonId uuid.UUID, skippedHours uint64, isLegit bool) *Skipped {
	return &Skipped{studentId: studentId, groupId: groupId, lessonId: lessonId, skippedHours: skippedHours}
}

func (sk *Skipped) StudentId() uuid.UUID {
	return sk.studentId
}
func (sk *Skipped) GroupId() uuid.UUID {
	return sk.groupId
}

func (sk *Skipped) LessonId() uuid.UUID {
	return sk.lessonId
}

func (sk *Skipped) IsLegit() bool {
	return sk.isLegit
}

func (sk *Skipped) SkippedHours() uint64 {
	return sk.skippedHours
}
