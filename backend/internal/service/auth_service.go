package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/repository"
)

// AuthService — регистрация, вход и управление JWT-токенами.
// Пароли хэшируются bcrypt до сохранения; открытые пароли никогда не сохраняются.
// JWT-токены используют HMAC-SHA256 с настраиваемым сроком действия (по умолчанию 24h).
type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

// AuthClaims — содержимое JWT с ID пользователя и ролью.
// Эти данные помещаются в контекст запроса Auth middleware.
type AuthClaims struct {
	UserID uuid.UUID      `json:"user_id"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthResponse возвращается при успешном входе, содержит JWT-токен
// и публичные данные пользователя (без хэша пароля).
type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

// Register создаёт нового пользователя с указанными учётными данными и ролью.
// Пароль хэшируется bcrypt. Возвращает созданного пользователя
// с очищенным PasswordHash (никогда не отдаётся в API).
func (s *AuthService) Register(ctx context.Context, email, password string, role model.UserRole, lastName, firstName string) (*model.User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		LastName:     lastName,
		FirstName:    firstName,
		IsActive:     true,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

// Login аутентифицирует пользователя по email и паролю.
// Возвращает JWT-токен (истекает через 24ч) и данные пользователя.
// Сообщения об ошибках одинаковы для неверного email и пароля ("invalid email or password"),
// чтобы предотвратить перебор пользователей. Отключённые учётные записи отклоняются отдельно.
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !checkPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	claims := &AuthClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return &AuthResponse{Token: tokenStr, User: user}, nil
}

// ValidateToken парсит и проверяет JWT-токен.
// Возвращает claims (user_id, role) при успехе или ошибку,
// если токен невалиден, истёк или был подделан.
func (s *AuthService) ValidateToken(tokenStr string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}