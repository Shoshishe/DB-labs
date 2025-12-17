package entities

import (
	"time"

	"github.com/google/uuid"
)

type OmissionType string

func (tp OmissionType) String() string {
	switch tp {
	case Sick:
		return "Sick"
	case Meeting:
		return "Meeting"
	default:
		panic("Invalid value used for omission type")
	}
}

const (
	Sick OmissionType = "Sick"
	Meeting OmissionType = "Meeting"
)

type Omission struct {
	studentId    uuid.UUID
	groupId      uuid.UUID
	info         string
	startTime    time.Time
	endTime      time.Time
	omissionType OmissionType
}

func NewOmission(studentId, groupId uuid.UUID, info string, startTime, endTime time.Time, omissionType OmissionType) *Omission {
	return &Omission{studentId: studentId, groupId: groupId, info: info, startTime: startTime, endTime: endTime, omissionType: omissionType}
}
