package repository

import (
	"context"
	usr "wch/services/users/pkg/model"

	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRepositoryPg struct {
	db *sql.DB
}

func NewUserRepositoryPg(db *sql.DB) *UserRepositoryPg {
	return &UserRepositoryPg{db}
}

func (r *UserRepositoryPg) CreateUser(ctx context.Context, u *usr.User) error {
	var id uuid.UUID
	if err := r.db.QueryRow(
		`INSERT INTO users (id, username, email, password_hash, 
		first_name, last_name, surname, avatar_url,
		department_id, created_at, updated_at) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		u.ID,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.Surname,
		u.AvatarURL,
		//u.Status,
		u.DepartmentID,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryPg) GetUserByID(ctx context.Context, id string) (*usr.User, error) {
	u := &usr.User{}
	if err := r.db.QueryRow(
		`SELECT id, username, email, password_hash, 
		first_name, last_name, surname, avatar_url, 
		status, department_id, created_at, updated_at
		FROM users WHERE id=$1`,
		id,
	).Scan(&u); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepositoryPg) GetUserByEmail(ctx context.Context, email string) (*usr.User, error) {
	u := &usr.User{}
	if err := r.db.QueryRow(
		`SELECT id, username, email, password_hash, 
		first_name, last_name, surname, avatar_url, 
		department_id
		FROM users WHERE email=$1`,
		email,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.Surname,
		&u.AvatarURL,
		&u.DepartmentID,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepositoryPg) GetUserByName(ctx context.Context, username string) (*usr.User, error) {
	return nil, nil
}

func (r *UserRepositoryPg) UpdateUser(ctx context.Context, u *usr.User) error {
	return nil
}

func (r *UserRepositoryPg) DeleteUser(ctx context.Context, id string) error {
	return nil
}
