package model

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	AvatarURL string    `db:"avatar_url" json:"avatar_url"`
	Status    string    `db:"status" json:"status"`
}
