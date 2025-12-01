package repository

import (
	"context"
	"wch/pkg/models"
	dep "wch/pkg/models"

	"database/sql"

	_ "github.com/lib/pq"
)

type DepartmentRepositoryPg struct {
	db *sql.DB
}

func NewDepartmentRepositoryPg(db *sql.DB) *DepartmentRepositoryPg {
	return &DepartmentRepositoryPg{db: db}
}

func (r *DepartmentRepositoryPg) CreateDepartment(ctx context.Context, d *dep.Department) error {
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

func (r *DepartmentRepositoryPg) GetDepartmentByID(ctx context.Context, id int) (*dep.Department, error) {
	d := &dep.Department{}
	if err := r.db.QueryRow(
		"SELECT id, name, rescription FROM departments WHERE id=$1",
		id,
	).Scan(&d.ID, &d.Name, &d.Description); err != nil {
		return nil, err
	}

	return d, nil
}

func (r *DepartmentRepositoryPg) GetDepartmentByName(ctx context.Context, name string) (*dep.Department, error) {
	d := &dep.Department{}
	if err := r.db.QueryRow(
		"SELECT id, name, rescription FROM departments WHERE name=$1",
		name,
	).Scan(&d.ID, &d.Name, &d.Description); err != nil {
		return nil, err
	}

	return d, nil
}

func (r *DepartmentRepositoryPg) UpdateDepartment(ctx context.Context, d *dep.Department) error {
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

func (r *DepartmentRepositoryPg) DeleteDepartment(ctx context.Context, id int) error {
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
func (r *DepartmentRepositoryPg) GetAll(ctx context.Context) ([]*models.Department, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM departments;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []*models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.Description); err != nil {
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
