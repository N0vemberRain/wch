package repository

import (
	"context"

	dep "wch/services/departments/pkg/model"
)

type Repository interface {
	CreateDepartment(ctx context.Context, d *dep.Department) error
	GetDepartmentByID(ctx context.Context, id int) (*dep.Department, error)
	GetDepartmentByName(ctx context.Context, name string) (*dep.Department, error)
	UpdateDepartment(ctx context.Context, d *dep.Department) error
	DeleteDepartment(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]*dep.Department, error)
}
