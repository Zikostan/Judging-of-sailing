// Package service — слой бизнес-логики приложения.
//
// Каждая доменная сущность имеет соответствующий сервис, который инкапсулирует
// бизнес-правила и делегирует доступ к данным репозиториям.
//
// Структура Services агрегирует все сервисы и предоставляет единую точку сборки.
// Используйте service.New(repos, jwtSecret) для связывания всех зависимостей.
//
// Типовые ошибки экспортируются как переменные пакета для единообразной
// обработки в слое хендлеров: ErrNotFound, ErrDuplicate, ErrInvalidInput, ErrForbidden.
package service

import (
	"errors"

	"github.com/user/judging-of-sailing/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound — запрашиваемый ресурс не существует.
	ErrNotFound = errors.New("not found")

	// ErrDuplicate — нарушение уникальности (дубликат email, участника в регате и т.д.).
	ErrDuplicate = errors.New("duplicate entry")

	// ErrInvalidInput — данные запроса не прошли валидацию (неверный пол, пропущены поля).
	ErrInvalidInput = errors.New("invalid input")

	// ErrForbidden — у пользователя нет прав на запрошенное действие.
	ErrForbidden = errors.New("forbidden")
)

// Services агрегирует все доменные сервисы в одну структуру.
// Создавайте через New() с заполненным Repositories и JWT-секретом.
type Services struct {
	Auth      *AuthService
	Regatta   *RegattaService
	Sailor    *SailorService
	Race      *RaceService
	RaceGroup *RaceGroupService
	Category  *CategoryService
	Document  *DocumentService
}

// New создаёт экземпляры всех сервисов, связывая репозитории и общие зависимости.
// jwtSecret передаётся в AuthService для подписи и проверки токенов.
func New(repos *repository.Repositories, jwtSecret string) *Services {
	return &Services{
		Auth:      NewAuthService(repos.User, jwtSecret),
		Regatta:   NewRegattaService(repos.Regatta),
		Sailor:    NewSailorService(repos.Sailor),
		Race:      NewRaceService(repos.Race, repos.RaceGroup),
		RaceGroup: NewRaceGroupService(repos.RaceGroup),
		Category:  NewCategoryService(repos.Category),
		Document:  NewDocumentService(repos.Document),
	}
}

// hashPassword возвращает bcrypt-хэш пароля.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword сравнивает пароль в открытом виде с bcrypt-хэшем. Возвращает true при совпадении.
func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// validateGender проверяет, что пол указан как "male" или "female".
func validateGender(gender string) bool {
	return gender == "male" || gender == "female"
}