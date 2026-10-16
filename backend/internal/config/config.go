// Package config — конфигурация приложения, загружаемая из переменных окружения.
//
// Секреты (DATABASE_URL, JWT_SECRET) намеренно НЕ имеют значений по умолчанию:
// если переменная не задана, приложение завершается с ошибкой. Это исключает
// случайный запуск в продакшене с дефолтным паролем БД или ключом подписи JWT.
//
// Для локальной разработки значения удобно держать в файле .env
// (шаблон — .env.example). Файл .env читается автоматически через godotenv
// из текущей директории и не должен попадать в систему контроля версий.
//
// Переменные окружения:
//   - SERVER_ADDR     — TCP-адрес HTTP-сервера (по умолчанию ":8080")
//   - DATABASE_URL    — строка подключения к PostgreSQL (обязательна)
//   - JWT_SECRET      — HMAC-ключ для подписи JWT (обязателен)
//   - JWT_EXPIRATION  — время жизни токена в формате Go duration (по умолчанию "24h")
//   - MIGRATION_PATH  — путь к файлам миграций SQL (по умолчанию "file://migrations")
//   - ALLOWED_ORIGINS — значение заголовка Access-Control-Allow-Origin (по умолчанию "*")
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	// Пустой импорт подключает autoload: при старте godotenv загружает .env
	// из текущей рабочей директории, если файл существует. Переменные,
	// уже заданные в окружении процесса, не перезаписываются.
	_ "github.com/joho/godotenv/autoload"
)

// Config содержит всю конфигурацию приложения.
type Config struct {
	// ServerAddr — TCP-адрес для HTTP-сервера, например ":8080".
	ServerAddr string

	// DatabaseURL — строка подключения к PostgreSQL. Обязателен.
	DatabaseURL string

	// JWTSecret — HMAC-SHA256 ключ для подписи и верификации JWT-токенов. Обязателен.
	JWTSecret string

	// JWTExpiration — время жизни JWT-токена после выпуска.
	JWTExpiration time.Duration

	// MigrationPath — путь к директории с SQL-миграциями.
	MigrationPath string

	// AllowedOrigins — значение заголовка Access-Control-Allow-Origin для CORS.
	AllowedOrigins string
}

// Load читает конфигурацию из переменных окружения и возвращает Config.
// Переменные без чувствительных данных заменяются значениями по умолчанию,
// а обязательные секреты (DATABASE_URL, JWT_SECRET) завершают процесс,
// если не заданы.
func Load() Config {
	return Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:    mustEnv("DATABASE_URL"),
		JWTSecret:      mustEnv("JWT_SECRET"),
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

// mustEnv возвращает значение переменной окружения или завершает процесс,
// если она не задана. Используется для обязательных секретов, у которых
// не должно быть небезопасного значения по умолчанию.
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
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
