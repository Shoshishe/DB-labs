package stored

import (
	"time"

	"github.com/google/uuid"
)

type Omission struct {
	StudentId    uuid.UUID
	GroupId      uuid.UUID
	Info         string
	StartTime    time.Time
	EndTime      time.Time
	Type   string
}
