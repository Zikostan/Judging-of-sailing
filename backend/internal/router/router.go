// Package router — связывает HTTP-маршруты с хендлерами и применяет middleware.
//
// Роутер использует встроенный в Go 1.22+ http.ServeMux с поддержкой
// паттернов (метод + путь). Внешние библиотеки роутинга не требуются.
//
// Структура маршрутов:
//   - Публичные (без аутентификации): POST /api/v1/auth/*
//   - Защищённые (JWT обязателен):        /api/v1/*
//   - Swagger UI:                         /swagger/*
//
// Цепочка middleware: CORS → внешний mux → Auth → защищённый mux → хендлер
//
// Swagger UI доступен по /swagger/ после запуска сервера.
package router

import (
	"net/http"

	"github.com/user/judging-of-sailing/backend/docs"
	"github.com/user/judging-of-sailing/backend/internal/handler"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

const swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Judging of Sailing API — Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({ url: "/swagger/doc.json", dom_id: "#swagger-ui" });
  </script>
</body>
</html>`

// New создаёт и возвращает полный HTTP-обработчик со всеми маршрутами и middleware.
// authMiddleware и corsMiddleware передаются из main.go, чтобы разорвать
// циклическую зависимость между пакетами router и middleware.
func New(
	svcs *service.Services,
	authMiddleware func(http.Handler) http.Handler,
	corsMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	authH := handler.NewAuthHandler(svcs.Auth)
	regattaH := handler.NewRegattaHandler(svcs.Regatta)
	sailorH := handler.NewSailorHandler(svcs.Sailor)
	raceGroupH := handler.NewRaceGroupHandler(svcs.RaceGroup)
	raceH := handler.NewRaceHandler(svcs.Race)
	categoryH := handler.NewCategoryHandler(svcs.Category)
	documentH := handler.NewDocumentHandler(svcs.Document)

	// Swagger
	spec := []byte(docs.SwaggerInfo.ReadDoc())
	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(spec)
	})
	mux.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(swaggerUI))
	})

	// Public routes
	mux.HandleFunc("POST /api/v1/auth/register", authH.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)

	// Protected routes
	protected := http.NewServeMux()

	// Regatta
	protected.HandleFunc("GET /api/v1/regattas", regattaH.List)
	protected.HandleFunc("POST /api/v1/regattas", regattaH.Create)
	protected.HandleFunc("GET /api/v1/regattas/{id}", regattaH.GetByID)
	protected.HandleFunc("PUT /api/v1/regattas/{id}", regattaH.Update)
	protected.HandleFunc("DELETE /api/v1/regattas/{id}", regattaH.Delete)

	// Category
	protected.HandleFunc("GET /api/v1/categories", categoryH.List)
	protected.HandleFunc("POST /api/v1/categories", categoryH.Create)
	protected.HandleFunc("GET /api/v1/categories/{id}", categoryH.GetByID)
	protected.HandleFunc("PUT /api/v1/categories/{id}", categoryH.Update)
	protected.HandleFunc("DELETE /api/v1/categories/{id}", categoryH.Delete)

	// Sailor
	protected.HandleFunc("GET /api/v1/regattas/{regatta_id}/sailors", sailorH.ListByRegatta)
	protected.HandleFunc("GET /api/v1/regattas/{regatta_id}/sailors/find", sailorH.FindByClassNumber)
	protected.HandleFunc("POST /api/v1/sailors", sailorH.Create)
	protected.HandleFunc("GET /api/v1/sailors/{id}", sailorH.GetByID)
	protected.HandleFunc("PUT /api/v1/sailors/{id}", sailorH.Update)
	protected.HandleFunc("DELETE /api/v1/sailors/{id}", sailorH.Delete)

	// Race Group
	protected.HandleFunc("GET /api/v1/regattas/{regatta_id}/race-groups", raceGroupH.ListByRegatta)
	protected.HandleFunc("POST /api/v1/race-groups", raceGroupH.Create)
	protected.HandleFunc("GET /api/v1/race-groups/{id}", raceGroupH.GetByID)
	protected.HandleFunc("POST /api/v1/race-groups/{id}/start", raceGroupH.Start)
	protected.HandleFunc("POST /api/v1/race-groups/{id}/finish", raceGroupH.Finish)
	protected.HandleFunc("DELETE /api/v1/race-groups/{id}", raceGroupH.Delete)

	// Race
	protected.HandleFunc("GET /api/v1/race-groups/{group_id}/races", raceH.ListByGroup)
	protected.HandleFunc("POST /api/v1/races", raceH.Record)
	protected.HandleFunc("GET /api/v1/races/{id}", raceH.GetByID)
	protected.HandleFunc("PATCH /api/v1/races/{id}/result", raceH.UpdateResult)
	protected.HandleFunc("POST /api/v1/races/{id}/link-sailor", raceH.LinkSailor)
	protected.HandleFunc("DELETE /api/v1/races/{id}", raceH.Delete)

	// Document
	protected.HandleFunc("GET /api/v1/sailors/{sailor_id}/documents", documentH.ListBySailor)
	protected.HandleFunc("POST /api/v1/documents", documentH.Upload)
	protected.HandleFunc("GET /api/v1/documents/{id}", documentH.GetByID)
	protected.HandleFunc("POST /api/v1/documents/{id}/approve", documentH.Approve)
	protected.HandleFunc("POST /api/v1/documents/{id}/reject", documentH.Reject)
	protected.HandleFunc("POST /api/v1/documents/{id}/request-revision", documentH.RequestRevision)
	protected.HandleFunc("DELETE /api/v1/documents/{id}", documentH.Delete)

	// Chain middleware: CORS -> Auth -> Protected routes
	mux.Handle("/api/v1/", authMiddleware(protected))

	return corsMiddleware(mux)
}