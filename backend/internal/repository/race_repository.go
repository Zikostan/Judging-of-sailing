package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// RaceRepository — CRUD и операции связывания для индивидуальных замеров.
// Уникальное ограничение на (group_id, boat_class, sail_number) гарантирует,
// что каждая лодка замеряется только один раз за заезд.
type RaceRepository struct {
	pool *pgxpool.Pool
}

func NewRaceRepository(pool *pgxpool.Pool) *RaceRepository {
	return &RaceRepository{pool: pool}
}

// Create вставляет новый замер. sailor_id может быть NULL, если участник
// ещё не опознан на момент замера.
func (r *RaceRepository) Create(ctx context.Context, race *model.Race) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO race (group_id, sailor_id, judge_id, boat_class, sail_number, finish_time, place, penalty)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, created_at, updated_at`,
		race.GroupID, race.SailorID, race.JudgeID, race.BoatClass, race.SailNumber, race.FinishTime, race.Place, race.Penalty,
	).Scan(&race.ID, &race.CreatedAt, &race.UpdatedAt)
}

// GetByID возвращает замер по UUID.
func (r *RaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Race, error) {
	race := &model.Race{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, group_id, sailor_id, judge_id, boat_class, sail_number, finish_time, place, penalty, created_at, updated_at
		 FROM race WHERE id = $1`, id,
	).Scan(&race.ID, &race.GroupID, &race.SailorID, &race.JudgeID, &race.BoatClass, &race.SailNumber, &race.FinishTime, &race.Place, &race.Penalty, &race.CreatedAt, &race.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return race, nil
}

// ListByGroup возвращает все замеры в заезде, отсортированные по классу лодки и времени финиша.
// Это основной запрос для отображения результатов по заезду.
func (r *RaceRepository) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]model.Race, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, group_id, sailor_id, judge_id, boat_class, sail_number, finish_time, place, penalty, created_at, updated_at
		 FROM race WHERE group_id = $1 ORDER BY boat_class, finish_time`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []model.Race
	for rows.Next() {
		var race model.Race
		if err := rows.Scan(&race.ID, &race.GroupID, &race.SailorID, &race.JudgeID, &race.BoatClass, &race.SailNumber, &race.FinishTime, &race.Place, &race.Penalty, &race.CreatedAt, &race.UpdatedAt); err != nil {
			return nil, err
		}
		races = append(races, race)
	}
	return races, nil
}

// LinkSailor привязывает ранее неизвестного участника к замеру.
// Используется, когда судья записал лодку по классу+номеру, но личность
// участника устанавливается позже секретарём.
func (r *RaceRepository) LinkSailor(ctx context.Context, raceID, sailorID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE race SET sailor_id = $1 WHERE id = $2`, sailorID, raceID)
	return err
}

// UpdateResult обновляет время финиша, место и штраф для замера.
// Каждое поле — указатель; nil-значения оставляют существующее значение без изменений.
func (r *RaceRepository) UpdateResult(ctx context.Context, raceID uuid.UUID, finishTime *float32, place *int16, penalty *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE race SET finish_time = $1, place = $2, penalty = $3 WHERE id = $4`,
		finishTime, place, penalty, raceID)
	return err
}

// Delete удаляет замер по ID.
func (r *RaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM race WHERE id = $1`, id)
	return err
}