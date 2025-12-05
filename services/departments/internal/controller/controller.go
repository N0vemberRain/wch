package controller

import (
	"context"
	"errors"
	"strings"

	"wch/services/departments/internal/repository"
	"wch/services/departments/pkg/model"
)

// Controller defines business logic for departments
type Controller struct {
	repo repository.Repository
}

func NewController(repo repository.Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) Create(ctx context.Context, d *model.Department) error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("department name is required")
	}

	if exist, _ := c.repo.GetDepartmentByName(ctx, d.Name); exist != nil {
		return errors.New("department already exists")
	}

	return c.repo.CreateDepartment(ctx, d)
}

func (c *Controller) GetByID(ctx context.Context, id int) (*model.Department, error) {
	return c.repo.GetDepartmentByID(ctx, id)
}

func (c *Controller) GetByName(ctx context.Context, name string) (*model.Department, error) {
	return c.repo.GetDepartmentByName(ctx, name)
}

func (c *Controller) GetAll(ctx context.Context) ([]*model.Department, error) {
	return c.repo.GetAll(ctx)
}
