package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// RegattaRepository — CRUD-операции для регат.
// Регаты — корневая сущность, объединяющая заезды и участников.
type RegattaRepository struct {
	pool *pgxpool.Pool
}

func NewRegattaRepository(pool *pgxpool.Pool) *RegattaRepository {
	return &RegattaRepository{pool: pool}
}

// Create вставляет новую регату и заполняет её ID и временные метки через RETURNING.
func (r *RegattaRepository) Create(ctx context.Context, reg *model.Regatta) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO regatta (title, start_date, end_date, location)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, updated_at`,
		reg.Title, reg.StartDate, reg.EndDate, reg.Location,
	).Scan(&reg.ID, &reg.CreatedAt, &reg.UpdatedAt)
}

// GetByID возвращает регату по UUID. Возвращает nil и ошибку, если не найдена.
func (r *RegattaRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Regatta, error) {
	reg := &model.Regatta{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, start_date, end_date, location, created_at, updated_at
		 FROM regatta WHERE id = $1`, id,
	).Scan(&reg.ID, &reg.Title, &reg.StartDate, &reg.EndDate, &reg.Location, &reg.CreatedAt, &reg.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

// List возвращает все регаты, отсортированные по start_date по убыванию (сначала новые).
func (r *RegattaRepository) List(ctx context.Context) ([]model.Regatta, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, start_date, end_date, location, created_at, updated_at
		 FROM regatta ORDER BY start_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regs []model.Regatta
	for rows.Next() {
		var reg model.Regatta
		if err := rows.Scan(&reg.ID, &reg.Title, &reg.StartDate, &reg.EndDate, &reg.Location, &reg.CreatedAt, &reg.UpdatedAt); err != nil {
			return nil, err
		}
		regs = append(regs, reg)
	}
	return regs, nil
}

// Update изменяет все изменяемые поля регаты по её ID.
func (r *RegattaRepository) Update(ctx context.Context, reg *model.Regatta) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE regatta SET title=$1, start_date=$2, end_date=$3, location=$4 WHERE id=$5`,
		reg.Title, reg.StartDate, reg.EndDate, reg.Location, reg.ID)
	return err
}

// Delete удаляет регату и все каскадно связанные данные (заезды, замеры, участники, документы).
func (r *RegattaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM regatta WHERE id = $1`, id)
	return err
}