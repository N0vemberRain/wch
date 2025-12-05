package repository

import (
	"context"
	"wch/services/departments/pkg/model"

	"database/sql"

	_ "github.com/lib/pq"
)

type RepositoryPg struct {
	db *sql.DB
}

func NewRepositoryPg(db *sql.DB) *RepositoryPg {
	return &RepositoryPg{db: db}
}

func (r *RepositoryPg) CreateDepartment(ctx context.Context, d *model.Department) error {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO departments (name, description) VALUES ($1, $2) RETURNING id",
		d.Name,
		d.Description,
	).Scan(&id); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryPg) GetDepartmentByID(ctx context.Context, id int) (*model.Department, error) {
	d := &model.Department{}
	if err := r.db.QueryRow(
		"SELECT id, name FROM departments WHERE id=$1",
		id,
	).Scan(&d.ID, &d.Name); err != nil {
		return nil, err
	}

	return d, nil
}

func (r *RepositoryPg) GetDepartmentByName(ctx context.Context, name string) (*model.Department, error) {
	d := &model.Department{}
	if err := r.db.QueryRow(
		"SELECT id, name FROM departments WHERE name=$1",
		name,
	).Scan(&d.ID, &d.Name); err != nil {
		return nil, err
	}

	return d, nil
}

func (r *RepositoryPg) UpdateDepartment(ctx context.Context, d *model.Department) error {
	_, err := r.db.Exec(
		"UPDATE departments SET name=$1, description=$2 WHERE id=$3",
		d.Name,
		d.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryPg) DeleteDepartment(ctx context.Context, id int) error {
	_, err := r.db.Exec(
		"DELETE FROM departments SET WHERE id=$3",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

// There might be an error
func (r *RepositoryPg) GetAll(ctx context.Context) ([]*model.Department, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM departments;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []*model.Department
	for rows.Next() {
		var d model.Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		deps = append(deps, &d)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deps, nil
}
