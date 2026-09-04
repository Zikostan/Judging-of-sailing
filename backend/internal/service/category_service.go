package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// CategoryService — управление зачётными категориями (деление по возрасту/полу/классу).
// Категории используются заездами для определения дивизиона.
type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, c *model.Category) error {
	return s.repo.Create(ctx, c)
}

func (s *CategoryService) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]model.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Update(ctx context.Context, c *model.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}