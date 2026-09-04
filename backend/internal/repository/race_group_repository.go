package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// RaceGroupRepository — CRUD и управление статусами заездов (heats).
// Жизненный цикл заезда: scheduled → running → finished (или cancelled).
type RaceGroupRepository struct {
	pool *pgxpool.Pool
}

func NewRaceGroupRepository(pool *pgxpool.Pool) *RaceGroupRepository {
	return &RaceGroupRepository{pool: pool}
}

// Create вставляет новый заезд. Номер должен быть уникальным в рамках регаты.
func (r *RaceGroupRepository) Create(ctx context.Context, rg *model.RaceGroup) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO race_group (regatta_id, category_id, number, start_at, status)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		rg.RegattaID, rg.CategoryID, rg.Number, rg.StartAt, rg.Status,
	).Scan(&rg.ID, &rg.CreatedAt, &rg.UpdatedAt)
}

// GetByID возвращает заезд по UUID.
func (r *RaceGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.RaceGroup, error) {
	rg := &model.RaceGroup{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, regatta_id, category_id, number, start_at, status, created_at, updated_at
		 FROM race_group WHERE id = $1`, id,
	).Scan(&rg.ID, &rg.RegattaID, &rg.CategoryID, &rg.Number, &rg.StartAt, &rg.Status, &rg.CreatedAt, &rg.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rg, nil
}

// ListByRegatta возвращает все заезды регаты, отсортированные по номеру.
func (r *RaceGroupRepository) ListByRegatta(ctx context.Context, regattaID uuid.UUID) ([]model.RaceGroup, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, regatta_id, category_id, number, start_at, status, created_at, updated_at
		 FROM race_group WHERE regatta_id = $1 ORDER BY number`, regattaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.RaceGroup
	for rows.Next() {
		var rg model.RaceGroup
		if err := rows.Scan(&rg.ID, &rg.RegattaID, &rg.CategoryID, &rg.Number, &rg.StartAt, &rg.Status, &rg.CreatedAt, &rg.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, rg)
	}
	return groups, nil
}

// UpdateStatus изменяет статус заезда (например, с scheduled на running).
func (r *RaceGroupRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.RaceGroupStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE race_group SET status = $1 WHERE id = $2`, status, id)
	return err
}

// Delete удаляет заезд и все каскадно связанные замеры.
func (r *RaceGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM race_group WHERE id = $1`, id)
	return err
}