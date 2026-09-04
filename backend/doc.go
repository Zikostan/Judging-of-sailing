// Package backend — серверная часть приложения для судейства парусных соревнований.
//
// # Быстрый старт
//
//	cd backend
//	export DATABASE_URL="postgres://postgres:postgres@localhost:5432/judging?sslmode=disable"
//	go run ./cmd/server
//
// После запуска:
//   - API:           http://localhost:8080/api/v1
//   - Swagger UI:    http://localhost:8080/swagger/index.html
//
// # Структура проекта
//
//	cmd/server/main.go          — точка входа, инициализация и запуск
//	internal/
//	  config/config.go           — конфигурация из env-переменных
//	  model/models.go            — доменные типы и бизнес-правила
//	  repository/                — слой доступа к данным (PostgreSQL + pgx)
//	    pool.go                  — пул соединений
//	    repositories.go          — агрегатор репозиториев
//	    *_repository.go          — CRUD для каждой сущности
//	  service/                   — слой бизнес-логики
//	    service.go               — агрегатор сервисов, общие ошибки, хелперы
//	    *_service.go             — логика для каждой сущности
//	  handler/                   — HTTP-хендлеры (JSON API)
//	    helpers.go               — writeJSON, writeError, decodeJSON
//	    *_handler.go             — эндпоинты для каждой сущности
//	  middleware/auth.go         — JWT-аутентификация + CORS
//	  migration/migration.go     — применение SQL-миграций при старте
//	  router/router.go           — маршрутизация (Go 1.22+ ServeMux)
//	migrations/                  — SQL-миграции (001_create_tables.*.sql)
//	docs/                        — сгенерированная Swagger-спецификация
//
// # Слои и зависимости
//
//	HTTP-запрос → router → middleware (CORS → Auth) → handler → service → repository → PostgreSQL
//
// # Доменная модель
//
//	Регата (Regatta) ─── Заезд (RaceGroup) ─── Замер (Race)
//	Регата ─── Участник (Sailor)
//	Участник ─── Документ (Document)
//	Категория (Category) ─── Заезд
//	Пользователь (User) ─── Замер (как судья)
//	Пользователь ─── Документ (как проверяющий)
//	Участник ───╼ Категория (M:N, через sailor_category)
//
// # Ключевые бизнес-правила
//
//   - Участник идентифицируется по (регата, класс_лодки, номер_на_парусе)
//   - Замер можно записать без привязки к участнику (sailor_id = NULL),
//     привязка выполняется позже через LinkSailor
//   - Один номер на парусе может быть в разных классах (Optimist #10 и Laser #10)
//   - Возраст хранится как дата рождения, а не число
//   - Документы проходят workflow: pending → approved / rejected / need_revision
//
// # API эндпоинты (все под /api/v1)
//
//	Публичные:
//	  POST /auth/register        — регистрация пользователя
//	  POST /auth/login           — вход, получение JWT
//
//	Защищённые (требуют Bearer-токен):
//	  CRUD /regattas             — регаты
//	  CRUD /categories           — категории
//	  CRUD /sailors              — участники
//	  GET  /regattas/:id/sailors — список участников регаты
//	  GET  /regattas/:id/sailors/find?boat_class=&sail_number= — поиск по классу+номеру
//	  CRUD /race-groups          — заезды
//	  POST /race-groups/:id/start      — начать заезд
//	  POST /race-groups/:id/finish     — завершить заезд
//	  CRUD /races                — замеры
//	  PATCH /races/:id/result         — обновить результат замера
//	  POST /races/:id/link-sailor     — привязать участника к замеру
//	  CRUD /documents            — документы
//	  POST /documents/:id/approve          — утвердить документ
//	  POST /documents/:id/reject           — отклонить документ
//	  POST /documents/:id/request-revision — запросить доработку
//
// # Просмотр GoDoc-документации
//
//	cd backend
//	go doc ./internal/model           # все типы
//	go doc ./internal/service AuthService  # сервис аутентификации
//	go doc ./internal/repository RaceRepository  # репозиторий замеров
//	go doc -all ./internal/middleware  # все функции middleware
//
// # Просмотр Swagger
//
//	Запустите сервер и откройте http://localhost:8080/swagger/index.html
package backend