package uuid

import (
	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func New() *UUID {
	return &UUID{uuid.New()}
}
