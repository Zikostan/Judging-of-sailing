package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// SailorRepository — CRUD и поисковые операции для участников.
// Ключевой бизнес-запрос — FindByClassAndNumber, используемый судьями
// при записи замеров для опознания участника по видимым маркировкам лодки.
type SailorRepository struct {
	pool *pgxpool.Pool
}

func NewSailorRepository(pool *pgxpool.Pool) *SailorRepository {
	return &SailorRepository{pool: pool}
}

// Create вставляет нового участника. Уникальное ограничение (regatta_id, boat_class, sail_number)
// предотвращает дублирование регистраций в рамках одной регаты.
func (r *SailorRepository) Create(ctx context.Context, s *model.Sailor) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO sailor (last_name, first_name, middle_name, gender, birth_date, boat_class, sail_number, regatta_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, created_at, updated_at`,
		s.LastName, s.FirstName, s.MiddleName, s.Gender, s.BirthDate, s.BoatClass, s.SailNumber, s.RegattaID,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

// GetByID возвращает участника по UUID.
func (r *SailorRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Sailor, error) {
	s := &model.Sailor{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, last_name, first_name, middle_name, gender, birth_date, boat_class, sail_number, regatta_id, created_at, updated_at
		 FROM sailor WHERE id = $1`, id,
	).Scan(&s.ID, &s.LastName, &s.FirstName, &s.MiddleName, &s.Gender, &s.BirthDate, &s.BoatClass, &s.SailNumber, &s.RegattaID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// FindByClassAndNumber ищет участника по классу лодки и номеру на парусе в рамках конкретной регаты.
// Это основной поиск для судей при записи замеров. Использует индекс (regatta_id, boat_class, sail_number).
func (r *SailorRepository) FindByClassAndNumber(ctx context.Context, regattaID uuid.UUID, boatClass, sailNumber string) (*model.Sailor, error) {
	s := &model.Sailor{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, last_name, first_name, middle_name, gender, birth_date, boat_class, sail_number, regatta_id, created_at, updated_at
		 FROM sailor WHERE regatta_id = $1 AND boat_class = $2 AND sail_number = $3`,
		regattaID, boatClass, sailNumber,
	).Scan(&s.ID, &s.LastName, &s.FirstName, &s.MiddleName, &s.Gender, &s.BirthDate, &s.BoatClass, &s.SailNumber, &s.RegattaID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ListByRegatta возвращает всех участников регаты, отсортированных по фамилии и имени.
func (r *SailorRepository) ListByRegatta(ctx context.Context, regattaID uuid.UUID) ([]model.Sailor, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, last_name, first_name, middle_name, gender, birth_date, boat_class, sail_number, regatta_id, created_at, updated_at
		 FROM sailor WHERE regatta_id = $1 ORDER BY last_name, first_name`, regattaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sailors []model.Sailor
	for rows.Next() {
		var s model.Sailor
		if err := rows.Scan(&s.ID, &s.LastName, &s.FirstName, &s.MiddleName, &s.Gender, &s.BirthDate, &s.BoatClass, &s.SailNumber, &s.RegattaID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sailors = append(sailors, s)
	}
	return sailors, nil
}

// Update изменяет все поля участника по ID.
func (r *SailorRepository) Update(ctx context.Context, s *model.Sailor) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sailor SET last_name=$1, first_name=$2, middle_name=$3, gender=$4, birth_date=$5, boat_class=$6, sail_number=$7 WHERE id=$8`,
		s.LastName, s.FirstName, s.MiddleName, s.Gender, s.BirthDate, s.BoatClass, s.SailNumber, s.ID)
	return err
}

// Delete удаляет участника и все каскадно связанные данные (документы, связи с замерами).
func (r *SailorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sailor WHERE id = $1`, id)
	return err
}