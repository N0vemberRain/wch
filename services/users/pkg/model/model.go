package model

import (
	"time"
	"wch/pkg/models"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID          `db:"id" json:"id"`
	Username     string             `db:"username" json:"username"`
	Email        string             `db:"email" json:"email"`
	PasswordHash string             `db:"password_hash" json:"-"`
	FirstName    string             `db:"first_name" json:"first_name"`
	LastName     string             `db:"last_name" json:"last_name"`
	Surname      string             `db:"surname" json:"surname"`
	AvatarURL    string             `db:"avatar_url" json:"avatar_url"`
	Status       string             `db:"status" json:"status"`
	CreatedAt    time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `db:"update_at" json:"update_at"`
	DepartmentID int                `db:"department_id" json:"department_id"`
	Department   *models.Department `json: "department,omitempty"`
}

func NewUser(username, email, passwdHash, firstName, lastName, surname string) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwdHash,
		FirstName:    firstName,
		LastName:     lastName,
		Status:       "offline",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
