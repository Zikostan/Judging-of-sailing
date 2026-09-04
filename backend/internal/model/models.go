// Package model — основные доменные типы, используемые во всём приложении.
//
// Связи между сущностями:
//
//	Регата  ──1:N── Заезд ──1:N── Замер
//	Регата  ──1:N── Участник
//	Участник ──1:N── Документ
//	Пользователь ──1:N── Замер (как судья)
//	Пользователь ──1:N── Документ (как проверяющий)
//	Категория ──1:N── Заезд
//	Участник ──M:N── Категория (через связку sailor_category)
//
// Ключевые бизнес-правила:
//   - Участник идентифицируется по (regatta_id, boat_class, sail_number) — уникально в рамках регаты.
//   - Замер идентифицируется по (group_id, boat_class, sail_number) — уникально в рамках заезда.
//   - Класс лодки и номер на парусе дублируются в Race, чтобы можно было
//     делать замеры без предварительной регистрации участника (sailor_id может быть NULL).
//   - Возраст хранится как дата рождения, а не число (возраст меняется со временем).
package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole — роль пользователя в системе.
//   - judge: записывает замеры, просматривает результаты
//   - secretary: управляет участниками и документами, проверяет бумаги
//   - admin: полный доступ ко всем ресурсам
type UserRole string

const (
	RoleJudge     UserRole = "judge"
	RoleSecretary UserRole = "secretary"
	RoleAdmin     UserRole = "admin"
)

// DocumentStatus — жизненный цикл документа участника.
//   - pending: документ загружен, ожидает проверки
//   - approved: документ проверен и принят
//   - rejected: документ отклонён, повторная загрузка невозможна
//   - need_revision: документ требует доработки, участник может перезагрузить
type DocumentStatus string

const (
	DocPending      DocumentStatus = "pending"
	DocApproved     DocumentStatus = "approved"
	DocRejected     DocumentStatus = "rejected"
	DocNeedRevision DocumentStatus = "need_revision"
)

// RaceGroupStatus — жизненный цикл заезда (heat).
//   - scheduled: заезд создан, но ещё не начат
//   - running: заезд идёт, можно записывать замеры
//   - finished: заезд завершён, новые замеры запрещены
//   - cancelled: заезд отменён, все связанные замеры недействительны
type RaceGroupStatus string

const (
	RGScheduled RaceGroupStatus = "scheduled"
	RGRunning   RaceGroupStatus = "running"
	RGFinished  RaceGroupStatus = "finished"
	RGCancelled RaceGroupStatus = "cancelled"
)

// User — системный пользователь (судья, секретарь или администратор).
// PasswordHash никогда не сериализуется в JSON (тег json:"-").
// IsActive управляет возможностью входа — отключённые учётные записи
// отклоняются при логине, а не удаляются.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role"`
	LastName     string    `json:"last_name"`
	FirstName    string    `json:"first_name"`
	MiddleName   *string   `json:"middle_name,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Regatta — соревнование (регата).
// Содержит несколько заездов и зарегистрированных участников.
// start_date должна быть не позже end_date (проверяется CHECK-ограничением в БД).
type Regatta struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Location  *string   `json:"location,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category — зачётная категория (дивизион) в рамках регаты.
// Каждая категория задаёт пол, возрастной диапазон и класс лодки.
// Участники привязываются к категориям через таблицу sailor_category.
// Комбинация (name, boat_class) уникальна, например "Юниоры" + "Optimist".
type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	AgeFrom   int       `json:"age_from"`
	AgeTo     int       `json:"age_to"`
	BoatClass string    `json:"boat_class"`
}

// Sailor — участник соревнований (спортсмен), зарегистрированный в регате.
// Каждый участник уникально идентифицируется по (regatta_id, boat_class, sail_number) —
// один и тот же номер может быть в разных классах (например, "Optimist #10" и "Laser #10"),
// но не может повторяться в одном классе в рамках одной регаты.
// Пол должен быть "male" или "female" (проверяется в сервисном слое и CHECK в БД).
type Sailor struct {
	ID         uuid.UUID `json:"id"`
	LastName   string    `json:"last_name"`
	FirstName  string    `json:"first_name"`
	MiddleName *string   `json:"middle_name,omitempty"`
	Gender     string    `json:"gender"`
	BirthDate  time.Time `json:"birth_date"`
	BoatClass  string    `json:"boat_class"`
	SailNumber string    `json:"sail_number"`
	RegattaID  uuid.UUID `json:"regatta_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Document — документ (файл), связанный с участником.
// Документы проходят workflow проверки: pending → approved/rejected/need_revision.
// checked_by ссылается на секретаря, проверившего документ.
// Файл хранится внешним образом; file_url содержит ссылку на него.
type Document struct {
	ID        uuid.UUID       `json:"id"`
	SailorID  uuid.UUID       `json:"sailor_id"`
	Type      string          `json:"type"`
	FileURL   string          `json:"file_url"`
	Status    DocumentStatus  `json:"status"`
	CheckedBy *uuid.UUID      `json:"checked_by,omitempty"`
	Comment   *string         `json:"comment,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// RaceGroup — заезд (heat). Группа лодок, стартующих одновременно.
// Все замеры внутри группы имеют одинаковое время старта и категорию.
// Номер заезда уникален в рамках регаты (например, Заезд 1, Заезд 2...).
// Статус переключается: scheduled → running → finished (или cancelled).
type RaceGroup struct {
	ID         uuid.UUID        `json:"id"`
	RegattaID  uuid.UUID        `json:"regatta_id"`
	CategoryID uuid.UUID        `json:"category_id"`
	Number     int              `json:"number"`
	StartAt    time.Time        `json:"start_at"`
	Status     RaceGroupStatus  `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// Race — индивидуальный замер одной лодки в рамках заезда.
// Каждый замер фиксирует результат одной лодки: время финиша, место и штраф.
// boat_class и sail_number дублируются из записи участника, чтобы можно было
// делать замеры, когда участник ещё не известен (sailor_id = NULL).
// judge_id автоматически проставляется из аутентифицированного пользователя.
// Комбинация (group_id, boat_class, sail_number) уникальна в рамках заезда.
type Race struct {
	ID         uuid.UUID  `json:"id"`
	GroupID    uuid.UUID  `json:"group_id"`
	SailorID   *uuid.UUID `json:"sailor_id,omitempty"`
	JudgeID    uuid.UUID  `json:"judge_id"`
	BoatClass  string     `json:"boat_class"`
	SailNumber string     `json:"sail_number"`
	FinishTime *float32   `json:"finish_time,omitempty"`
	Place      *int16     `json:"place,omitempty"`
	Penalty    *string    `json:"penalty,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// SailorCategory — связка участника и категории (отношение M:N).
// Один участник может быть в нескольких категориях (например, "Юниоры" и "Открытый класс"),
// и одна категория содержит много участников.
type SailorCategory struct {
	SailorID   uuid.UUID `json:"sailor_id"`
	CategoryID uuid.UUID `json:"category_id"`
}