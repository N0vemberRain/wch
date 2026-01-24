package postgres

import (
	"context"
	"database/sql"
	"wch/services/auth/internal/domain"
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
	*domain.Credentials, error,
) {
	query := `SELECT * FROM credentials WHERE email=$1;`
	row := r.db.QueryRowContext(ctx, query, email)
	cred := &domain.Credentials{}
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
