package model

type TaskStatus string

const (
	StatusPool  TaskStatus = "pool"
	StatusToday TaskStatus = "today"
	StatusDone  TaskStatus = "done"
)

func (s TaskStatus) Valid() bool {
	switch s {
	case StatusPool, StatusToday, StatusDone:
		return true
	default:
		return false
	}
}

func (s TaskStatus) String() string {
	return string(s)
}
