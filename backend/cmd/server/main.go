// Package main — точка входа в приложение.
//
// Выполняет инициализацию подключения к БД, запуск миграций,
// сборку зависимостей: репозитории → сервисы → хендлеры → роутер,
// применение middleware и запуск HTTP-сервера с graceful shutdown.
//
// Конфигурация загружается из переменных окружения (см. internal/config).
// Сервер слушает на адресе из SERVER_ADDR (по умолчанию :8080).
package main

// @title           Judging of Sailing API
// @version         1.0.0
// @description     Backend API for the sailing competition judging application. Handles regattas, race groups, individual race measurements, sailors, documents, and user authentication.
// @termsOfService  https://github.com/user/judging-of-sailing

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                JWT token obtained from POST /api/v1/auth/login. Use format: Bearer <token>

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/judging-of-sailing/backend/internal/config"
	"github.com/user/judging-of-sailing/backend/internal/middleware"
	"github.com/user/judging-of-sailing/backend/internal/migration"
	"github.com/user/judging-of-sailing/backend/internal/repository"
	"github.com/user/judging-of-sailing/backend/internal/router"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("database connected")

	migrator := migration.NewRunner(pool)
	if err := migrator.Up(context.Background(), "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	repos := repository.New(pool)
	svcs := service.New(repos, cfg.JWTSecret)

	authMiddleware := middleware.Auth(svcs.Auth)
	corsMiddleware := middleware.CORS(cfg.AllowedOrigins)

	handler := router.New(svcs, authMiddleware, corsMiddleware)

	server := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on %s", cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-shutdown
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}