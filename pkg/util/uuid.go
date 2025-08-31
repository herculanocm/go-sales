package util

import "github.com/google/uuid"

type UUID = uuid.UUID

func New() UUID {
	return UUID(uuid.New())
}

func NewPtr() *UUID {
	id := New()
	return &id
}

func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID(id), nil
}
