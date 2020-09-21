package impl

import (
	"github.com/google/uuid"
)

func CreateStringContext(s string) func() string {
	return func() string {
		return s
	}
}

func CreateConnectionId() func() string {
	return func() string {
		return uuid.New().String()
	}
}
