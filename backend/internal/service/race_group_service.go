package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// RaceGroupService — управление жизненным циклом заездов (heats).
// Предоставляет явные методы Start и Finish для переключения статусов.
type RaceGroupService struct {
	repo *repository.RaceGroupRepository
}

func NewRaceGroupService(repo *repository.RaceGroupRepository) *RaceGroupService {
	return &RaceGroupService{repo: repo}
}

func (s *RaceGroupService) Create(ctx context.Context, rg *model.RaceGroup) error {
	return s.repo.Create(ctx, rg)
}

func (s *RaceGroupService) GetByID(ctx context.Context, id uuid.UUID) (*model.RaceGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RaceGroupService) ListByRegatta(ctx context.Context, regattaID uuid.UUID) ([]model.RaceGroup, error) {
	return s.repo.ListByRegatta(ctx, regattaID)
}

// StartRaceGroup переводит статус заезда в "running", разрешая запись замеров.
func (s *RaceGroupService) StartRaceGroup(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, id, model.RGRunning)
}

// FinishRaceGroup переводит статус заезда в "finished", запрещая новые замеры.
func (s *RaceGroupService) FinishRaceGroup(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, id, model.RGFinished)
}

func (s *RaceGroupService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}