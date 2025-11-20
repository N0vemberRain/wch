package controller

import (
	"context"
	"errors"
	"strings"

	"wch/pkg/models"
	"wch/services/users/internal/repository"
)

// DepartmentController defines business logic for departments
type DepartmentController struct {
	repo repository.DepartmentRepository
}

func NewDepartmentController(repo repository.DepartmentRepository) *DepartmentController {
	return &DepartmentController{repo: repo}
}

func (c *DepartmentController) Create(ctx context.Context, d *models.Department) error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("department name is required")
	}

	if exist, _ := c.repo.GetDepartmentByName(ctx, d.Name); exist != nil {
		return errors.New("department already exists")
	}

	return c.repo.CreateDepartment(ctx, d)
}

func (c *DepartmentController) GetByID(ctx context.Context, id int) (*models.Department, error) {
	return c.repo.GetDepartmentByID(ctx, id)
}

func (c *DepartmentController) GetByName(ctx context.Context, name string) (*models.Department, error) {
	return c.repo.GetDepartmentByName(ctx, name)
}

func (c *DepartmentController) GetAll(ctx context.Context) ([]*models.Department, error) {
	return c.repo.GetAll(ctx)
}
