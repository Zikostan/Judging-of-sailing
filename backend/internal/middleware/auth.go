// Package middleware — HTTP-обработчики для аутентификации, авторизации и CORS.
//
// Auth middleware извлекает JWT из заголовка Authorization (схема Bearer),
// проверяет его и добавляет user_id и user_role в контекст запроса.
// Все защищённые маршруты используют этот middleware; публичные маршруты
// (auth/register, auth/login) обходят его.
//
// CORS middleware устанавливает разрешающие кросс-доменные заголовки
// для доступа с фронтенда.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type contextKey string

const (
	// KeyUserID — ключ контекста для UUID аутентифицированного пользователя.
	KeyUserID contextKey = "user_id"

	// KeyUserRole — ключ контекста для роли пользователя (judge, secretary, admin).
	KeyUserRole contextKey = "user_role"
)

// Auth возвращает middleware, проверяющее Bearer JWT-токен из заголовка Authorization.
// При успехе хендлер может получить ID пользователя и роль через GetUserID и GetUserRole.
// При ошибке (нет заголовка, невалидный токен, истёкший токен) возвращает 401.
func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if tokenStr == header {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := authService.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), KeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, KeyUserRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID извлекает UUID аутентифицированного пользователя из контекста запроса.
// Возвращает UUID и true, если он присутствует, или нулевое значение и false.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(KeyUserID).(uuid.UUID)
	return id, ok
}

// GetUserRole извлекает роль пользователя из контекста запроса.
// Возвращает строку роли и true, или пустую строку и false.
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(KeyUserRole).(string)
	return role, ok
}

// CORS возвращает middleware, устанавливающее заголовки кросс-доменного доступа.
// Параметр allowedOrigins задаёт значение заголовка Access-Control-Allow-Origin.
// Preflight OPTIONS-запросы обрабатываются автоматически с ответом 204 No Content.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}