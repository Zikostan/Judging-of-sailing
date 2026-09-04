// Package repository — слой доступа к данным PostgreSQL для всех доменных сущностей.
//
// Каждая сущность имеет отдельный репозиторий с CRUD-операциями и
// доменными запросами (например, FindByClassAndNumber для участников).
// Все репозитории используют pgx/v5/pgxpool для пула соединений.
//
// Структура Repositories агрегирует все репозитории и является
// единой точкой создания для сервисного слоя.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool создаёт pgxpool.Pool из URL базы данных и проверяет соединение.
// MaxConns установлен в 20 — измените через конфигурацию pgxpool при необходимости.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 20
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}