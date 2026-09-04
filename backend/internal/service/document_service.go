package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// DocumentService — workflow проверки документов участников.
// Документы можно загружать, утверждать, отклонять или отправлять на доработку.
// Каждое изменение статуса фиксирует ID проверяющего секретаря.
type DocumentService struct {
	repo *repository.DocumentRepository
}

func NewDocumentService(repo *repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

// Upload создаёт запись нового документа. Файл должен храниться внешним образом.
func (s *DocumentService) Upload(ctx context.Context, doc *model.Document) error {
	return s.repo.Create(ctx, doc)
}

func (s *DocumentService) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DocumentService) ListBySailor(ctx context.Context, sailorID uuid.UUID) ([]model.Document, error) {
	return s.repo.ListBySailor(ctx, sailorID)
}

// Approve устанавливает статус документа "approved". Комментарий не требуется.
func (s *DocumentService) Approve(ctx context.Context, id uuid.UUID, checkedBy uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, id, model.DocApproved, checkedBy, nil)
}

// Reject устанавливает статус документа "rejected" с указанием причины.
func (s *DocumentService) Reject(ctx context.Context, id uuid.UUID, checkedBy uuid.UUID, comment string) error {
	return s.repo.UpdateStatus(ctx, id, model.DocRejected, checkedBy, &comment)
}

// RequestRevision устанавливает статус "need_revision" с инструкциями.
// После запроса доработки участник может загрузить исправленную версию.
func (s *DocumentService) RequestRevision(ctx context.Context, id uuid.UUID, checkedBy uuid.UUID, comment string) error {
	return s.repo.UpdateStatus(ctx, id, model.DocNeedRevision, checkedBy, &comment)
}

func (s *DocumentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}