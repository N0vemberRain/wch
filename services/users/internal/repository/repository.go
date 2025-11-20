package repository

import (
	"context"

	"wch/pkg/models"
	dep "wch/pkg/models"
	"wch/services/users/pkg/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByName(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, u *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, d *dep.Department) error
	GetDepartmentByID(ctx context.Context, id int) (*dep.Department, error)
	GetDepartmentByName(ctx context.Context, name string) (*dep.Department, error)
	UpdateDepartment(ctx context.Context, d *dep.Department) error
	DeleteDepartment(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]*models.Department, error)
}
