package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// SailorService — регистрация участников и поиск.
// Проверяет пол при Create и Update (должен быть "male" или "female").
// FindByClassAndNumber — ключевой метод, используемый судьями при замерах.
type SailorService struct {
	repo *repository.SailorRepository
}

func NewSailorService(repo *repository.SailorRepository) *SailorService {
	return &SailorService{repo: repo}
}

// Create регистрирует нового участника после проверки пола.
func (s *SailorService) Create(ctx context.Context, sailor *model.Sailor) error {
	if !validateGender(sailor.Gender) {
		return ErrInvalidInput
	}
	return s.repo.Create(ctx, sailor)
}

func (s *SailorService) GetByID(ctx context.Context, id uuid.UUID) (*model.Sailor, error) {
	return s.repo.GetByID(ctx, id)
}

// FindByClassAndNumber ищет участника по классу лодки и номеру на парусе в рамках регаты.
// Используется судьями, которые видят маркировку лодки, но не знают имени участника.
func (s *SailorService) FindByClassAndNumber(ctx context.Context, regattaID uuid.UUID, boatClass, sailNumber string) (*model.Sailor, error) {
	return s.repo.FindByClassAndNumber(ctx, regattaID, boatClass, sailNumber)
}

func (s *SailorService) ListByRegatta(ctx context.Context, regattaID uuid.UUID) ([]model.Sailor, error) {
	return s.repo.ListByRegatta(ctx, regattaID)
}

// Update изменяет данные участника после проверки пола.
func (s *SailorService) Update(ctx context.Context, sailor *model.Sailor) error {
	if !validateGender(sailor.Gender) {
		return ErrInvalidInput
	}
	return s.repo.Update(ctx, sailor)
}

func (s *SailorService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}