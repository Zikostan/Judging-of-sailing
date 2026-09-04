package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// RegattaService — бизнес-логика управления регатами.
// Пока является тонкой обёрткой над репозиторием; добавляйте валидацию
// и проверку прав здесь по мере роста приложения.
type RegattaService struct {
	repo *repository.RegattaRepository
}

func NewRegattaService(repo *repository.RegattaRepository) *RegattaService {
	return &RegattaService{repo: repo}
}

func (s *RegattaService) Create(ctx context.Context, reg *model.Regatta) error {
	return s.repo.Create(ctx, reg)
}

func (s *RegattaService) GetByID(ctx context.Context, id uuid.UUID) (*model.Regatta, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RegattaService) List(ctx context.Context) ([]model.Regatta, error) {
	return s.repo.List(ctx)
}

func (s *RegattaService) Update(ctx context.Context, reg *model.Regatta) error {
	return s.repo.Update(ctx, reg)
}

func (s *RegattaService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}