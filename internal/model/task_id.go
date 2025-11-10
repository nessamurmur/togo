package model

import (
	"github.com/google/uuid"
)

type TaskID uuid.UUID

func (t TaskID) String() string {
	return uuid.UUID(t).String()
}

func NewTaskID() TaskID {
	return TaskID(uuid.New())
}

func ParseTaskID(id string) (TaskID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return TaskID(uuid.Nil), err
	}
	return TaskID(uid), nil
}

func (t TaskID) IsEmpty() bool {
	return t == TaskID(uuid.Nil)
}

func (t TaskID) Equals(other TaskID) bool {
	return t == other
}

func (t TaskID) NotEquals(other TaskID) bool {
	return t != other
}
