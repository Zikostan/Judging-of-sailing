package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// DocumentRepository — CRUD и workflow-операции для документов участников.
// Документы проходят workflow проверки через UpdateStatus.
// checked_by фиксирует, какой секретарь выполнил проверку.
type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{pool: pool}
}

// Create вставляет новую запись документа. Файл должен храниться внешним образом;
// file_url — ссылка на внешнее хранилище.
func (r *DocumentRepository) Create(ctx context.Context, d *model.Document) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO document (sailor_id, type, file_url, status, checked_by, comment)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at, updated_at`,
		d.SailorID, d.Type, d.FileURL, d.Status, d.CheckedBy, d.Comment,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

// GetByID возвращает документ по UUID.
func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	d := &model.Document{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, sailor_id, type, file_url, status, checked_by, comment, created_at, updated_at
		 FROM document WHERE id = $1`, id,
	).Scan(&d.ID, &d.SailorID, &d.Type, &d.FileURL, &d.Status, &d.CheckedBy, &d.Comment, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// ListBySailor возвращает все документы участника, отсортированные по дате создания
// по убыванию (сначала новые).
func (r *DocumentRepository) ListBySailor(ctx context.Context, sailorID uuid.UUID) ([]model.Document, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, sailor_id, type, file_url, status, checked_by, comment, created_at, updated_at
		 FROM document WHERE sailor_id = $1 ORDER BY created_at DESC`, sailorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		var d model.Document
		if err := rows.Scan(&d.ID, &d.SailorID, &d.Type, &d.FileURL, &d.Status, &d.CheckedBy, &d.Comment, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// UpdateStatus изменяет статус документа и фиксирует проверяющего.
// Статус: DocApproved, DocRejected или DocNeedRevision.
// nil-комментарий допустим для утверждения; для отклонения и запроса правки требуется комментарий.
func (r *DocumentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.DocumentStatus, checkedBy uuid.UUID, comment *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE document SET status = $1, checked_by = $2, comment = $3 WHERE id = $4`,
		status, checkedBy, comment, id)
	return err
}

// Delete удаляет запись документа по ID.
func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM document WHERE id = $1`, id)
	return err
}