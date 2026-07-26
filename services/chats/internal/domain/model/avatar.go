package model

import "github.com/google/uuid"

type Avatar struct {
	Data     []byte
	MimeType string
	OwnerID  uuid.UUID
}
