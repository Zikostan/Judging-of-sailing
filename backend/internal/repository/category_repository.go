package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// CategoryRepository — CRUD-операции для зачётных категорий.
// Категории определяют деление на дивизионы по возрасту/полу/классу лодки.
type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

// Create вставляет новую категорию. Комбинация (name, boat_class) должна быть уникальна.
func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO category (name, gender, age_from, age_to, boat_class)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		c.Name, c.Gender, c.AgeFrom, c.AgeTo, c.BoatClass,
	).Scan(&c.ID)
}

// GetByID возвращает категорию по UUID.
func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	c := &model.Category{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, gender, age_from, age_to, boat_class FROM category WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Gender, &c.AgeFrom, &c.AgeTo, &c.BoatClass)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// List возвращает все категории, отсортированные по имени.
func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, gender, age_from, age_to, boat_class FROM category ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Gender, &c.AgeFrom, &c.AgeTo, &c.BoatClass); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

// Update изменяет все поля категории по ID.
func (r *CategoryRepository) Update(ctx context.Context, c *model.Category) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE category SET name=$1, gender=$2, age_from=$3, age_to=$4, boat_class=$5 WHERE id=$6`,
		c.Name, c.Gender, c.AgeFrom, c.AgeTo, c.BoatClass, c.ID)
	return err
}

// Delete удаляет категорию по ID. Заезды, ссылающиеся на эту категорию, должны быть обновлены заранее.
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM category WHERE id = $1`, id)
	return err
}