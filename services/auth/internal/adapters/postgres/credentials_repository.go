package postgres

import (
	"context"
	"database/sql"
	"wch/services/auth/internal/domain"
	"wch/services/auth/internal/domain/model"

	"github.com/google/uuid"
)

type CredentialsRepository struct {
	db *sql.DB
}

func NewCredentialsRepository(db *sql.DB) *CredentialsRepository {
	return &CredentialsRepository{
		db: db,
	}
}

func (r *CredentialsRepository) GetByEmail(ctx context.Context, email string) (
	*model.Credentials, error,
) {
	query := `SELECT user_id, email, password_hash, is_active, created_at FROM auth_credentials WHERE email=$1;`
	row := r.db.QueryRowContext(ctx, query, email)
	cred := &model.Credentials{}
	err := row.Scan(
		&cred.UserID,
		&cred.Name,
		&cred.PasswordHash,
		&cred.IsActive,
		&cred.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCredentialsNotFound
		} else {
			return nil, err
		}
	}

	return cred, nil
}

func (r *CredentialsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (
	*model.Credentials, error,
) {
	query := `SELECT user_id, email, password_hash, is_active, created_at FROM auth_credentials WHERE user_id=$1;`
	row := r.db.QueryRowContext(ctx, query, userID.String())
	cred := &model.Credentials{}
	err := row.Scan(
		&cred.UserID,
		&cred.Name,
		&cred.PasswordHash,
		&cred.IsActive,
		&cred.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCredentialsNotFound
		} else {
			return nil, err
		}
	}

	return cred, nil
}

func (r *CredentialsRepository) Save(ctx context.Context, cred *model.Credentials) error {
	query := `INSERT INTO auth_credentials (user_id, name, password_hash, created_at)
		VALUES ($1, $2, $3, $4);`
	return r.db.QueryRowContext(ctx, query, cred.UserID).Err()
}
