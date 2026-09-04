package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// RaceService — запись и управление индивидуальными замерами.
// Каждый замер фиксирует результат одной лодки в заезде.
// Сервис связывает судей (из контекста аутентификации) и участников (определяются пост-фактум).
type RaceService struct {
	raceRepo      *repository.RaceRepository
	raceGroupRepo *repository.RaceGroupRepository
}

func NewRaceService(raceRepo *repository.RaceRepository, raceGroupRepo *repository.RaceGroupRepository) *RaceService {
	return &RaceService{raceRepo: raceRepo, raceGroupRepo: raceGroupRepo}
}

// Record создаёт новый замер. judge_id должен быть установлен из аутентифицированного
// пользователя перед вызовом. sailor_id может быть nil.
func (s *RaceService) Record(ctx context.Context, race *model.Race) error {
	return s.raceRepo.Create(ctx, race)
}

func (s *RaceService) GetByID(ctx context.Context, id uuid.UUID) (*model.Race, error) {
	return s.raceRepo.GetByID(ctx, id)
}

func (s *RaceService) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]model.Race, error) {
	return s.raceRepo.ListByGroup(ctx, groupID)
}

// UpdateResult устанавливает время финиша, место и штраф для замера.
// Каждое поле — указатель; nil-значения оставляют существующее значение без изменений.
func (s *RaceService) UpdateResult(ctx context.Context, raceID uuid.UUID, finishTime *float32, place *int16, penalty *string) error {
	return s.raceRepo.UpdateResult(ctx, raceID, finishTime, place, penalty)
}

// LinkSailor привязывает ранее неизвестного участника к замеру.
// Вызывается секретарём после того, как личность участника установлена.
func (s *RaceService) LinkSailor(ctx context.Context, raceID, sailorID uuid.UUID) error {
	return s.raceRepo.LinkSailor(ctx, raceID, sailorID)
}

func (s *RaceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.raceRepo.Delete(ctx, id)
}