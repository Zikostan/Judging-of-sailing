// Package config — конфигурация приложения, загружаемая из переменных окружения.
//
// Все значения имеют разумные умолчания для локальной разработки.
// Переопределите любую переменную, установив соответствующий env-параметр.
//
// Переменные окружения:
//   - SERVER_ADDR     — TCP-адрес HTTP-сервера (по умолчанию ":8080")
//   - DATABASE_URL    — строка подключения к PostgreSQL (по умолчанию localhost:5432/judging)
//   - JWT_SECRET      — HMAC-ключ для подписи JWT (по умолчанию "change-me-in-production")
//   - JWT_EXPIRATION  — время жизни токена в формате Go duration (по умолчанию "24h")
//   - MIGRATION_PATH  — путь к файлам миграций SQL (по умолчанию "file://migrations")
//   - ALLOWED_ORIGINS — значение заголовка Access-Control-Allow-Origin (по умолчанию "*")
package config

import (
	"os"
	"strconv"
	"time"
)

// Config содержит всю конфигурацию приложения.
type Config struct {
	// ServerAddr — TCP-адрес для HTTP-сервера, например ":8080".
	ServerAddr string

	// DatabaseURL — строка подключения к PostgreSQL.
	DatabaseURL string

	// JWTSecret — HMAC-SHA256 ключ для подписи и верификации JWT-токенов.
	JWTSecret string

	// JWTExpiration — время жизни JWT-токена после выпуска.
	JWTExpiration time.Duration

	// MigrationPath — путь к директории с SQL-миграциями.
	MigrationPath string

	// AllowedOrigins — значение заголовка Access-Control-Allow-Origin для CORS.
	AllowedOrigins string
}

// Load читает конфигурацию из переменных окружения и возвращает Config.
// Отсутствующие переменные заменяются значениями по умолчанию.
func Load() Config {
	return Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:1AFOyzgR@localhost:5432/judje?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiration:  getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		MigrationPath:  getEnv("MIGRATION_PATH", "file://migrations"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
