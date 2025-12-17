package entities

type LessonType int8

const (
	Lecture LessonType = iota + 1
	Labwork
	Seminary
)

func (tp LessonType) String() string {
	switch tp {
	case Lecture:
		return "Lecture"
	case Labwork:
		return "Labwork"
	case Seminary:
		return "Seminary"
	default:
		panic("Invalid value used for omission type")
	}
}

type Lesson struct {
	lessonType LessonType
}

