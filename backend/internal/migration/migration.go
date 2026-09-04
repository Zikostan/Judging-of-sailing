// Package migration — применение SQL-миграций при запуске приложения.
//
// Миграции отслеживаются в таблице schema_migrations, которая создаётся
// автоматически перед выполнением миграций. Каждый .up.sql файл выполняется
// только один раз — повторный запуск сервера не вызывает ошибок.
//
// Соглашение об именах: NNN_description.up.sql (например, 001_create_tables.up.sql).
// Файлы выполняются в лексикографическом порядке.
package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Runner применяет SQL-миграции к базе PostgreSQL.
type Runner struct {
	pool *pgxpool.Pool
}

// NewRunner создаёт Runner, выполняющий миграции через указанный пул соединений.
func NewRunner(pool *pgxpool.Pool) *Runner {
	return &Runner{pool: pool}
}

// Up читает все *.up.sql файлы из migrationsDir, сортирует их по алфавиту
// и выполняет только те, которые ещё не были применены.
func (r *Runner) Up(ctx context.Context, migrationsDir string) error {
	if err := r.ensureTrackingTable(ctx); err != nil {
		return fmt.Errorf("ensure tracking table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, f := range upFiles {
		applied, err := r.isApplied(ctx, f)
		if err != nil {
			return fmt.Errorf("check applied %s: %w", f, err)
		}
		if applied {
			fmt.Printf("skipping already applied: %s\n", f)
			continue
		}

		path := filepath.Join(migrationsDir, f)
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}

		tx, err := r.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", f, err)
		}

		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("exec %s: %w", f, err)
		}

		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, f); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("track %s: %w", f, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit %s: %w", f, err)
		}

		fmt.Printf("applied migration: %s\n", f)
	}
	return nil
}

func (r *Runner) ensureTrackingTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename    TEXT PRIMARY KEY,
			applied_at  TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	return err
}

func (r *Runner) isApplied(ctx context.Context, filename string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`, filename,
	).Scan(&exists)
	return exists, err
}