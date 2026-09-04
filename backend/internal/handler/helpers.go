// Package handler — HTTP-обработчики для всех эндпоинтов API.
//
// Каждая доменная сущность имеет отдельный хендлер, связывающий HTTP-запросы
// с соответствующим сервисом. Хендлеры отвечают за:
//   - Разбор и валидацию тела запроса (JSON)
//   - Извлечение path/query параметров
//   - Вызов методов сервиса
//   - Сериализацию ответа (JSON) или сообщений об ошибках
//
// Типы запросов, специфичные для хендлеров (не входящие в model), определены
// в соответствующих файлах для ясности.
//
// Вспомогательные функции (writeJSON, writeError, decodeJSON) находятся в helpers.go.
package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON сериализует данные в JSON и отправляет с указанным HTTP-статусом.
// Устанавливает Content-Type: application/json; charset=utf-8.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError отправляет JSON-ответ с ошибкой и указанным HTTP-статусом.
// Формат: {"error": "<message>"}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// decodeJSON читает тело запроса и десериализует его в переданное значение.
// Тело запроса закрывается после чтения.
func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}