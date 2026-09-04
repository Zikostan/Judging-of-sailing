package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/judging-of-sailing/backend/internal/model"
)

// UserRepository — CRUD и аутентификационные запросы для системных пользователей.
// GetByEmail — основной поиск, используемый при входе в систему.
// password_hash никогда не раскрывается за пределами этого слоя.
type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create вставляет нового пользователя. Email должен быть уникальным (ограничение UNIQUE).
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO "user" (email, password_hash, role, last_name, first_name, middle_name, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		u.Email, u.PasswordHash, u.Role, u.LastName, u.FirstName, u.MiddleName, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// GetByID возвращает пользователя по UUID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, role, last_name, first_name, middle_name, is_active, created_at, updated_at
		 FROM "user" WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.LastName, &u.FirstName, &u.MiddleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByEmail находит пользователя по email. Используется сервисом аутентификации для входа.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, role, last_name, first_name, middle_name, is_active, created_at, updated_at
		 FROM "user" WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.LastName, &u.FirstName, &u.MiddleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// List возвращает всех пользователей, отсортированных по фамилии и имени.
func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, email, password_hash, role, last_name, first_name, middle_name, is_active, created_at, updated_at
		 FROM "user" ORDER BY last_name, first_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.LastName, &u.FirstName, &u.MiddleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}