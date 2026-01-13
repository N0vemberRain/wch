package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Cursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

func NewCursor(str string) (*Cursor, error) {
	if str == "" {
		return nil, errors.New("cursor is empty")
	}
	parts := strings.Split(str, "|")
	if len(parts) != 2 {
		return nil, errors.New("cursor is invalid")
	}
	tm, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, err
	}

	msg_id, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, err
	}

	return &Cursor{
		CreatedAt: tm,
		ID:        msg_id,
	}, nil
}
