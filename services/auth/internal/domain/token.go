package domain

import "time"

type Token struct {
	Value     string
	ExpiresAt time.Time
}
