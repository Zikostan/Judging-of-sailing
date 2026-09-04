package repository

import "github.com/jackc/pgx/v5/pgxpool"

// Repositories агрегирует все доменные репозитории в одну структуру.
// Используйте New() для создания всех репозиториев из одного пула соединений.
type Repositories struct {
	User      *UserRepository
	Regatta   *RegattaRepository
	Sailor    *SailorRepository
	Race      *RaceRepository
	RaceGroup *RaceGroupRepository
	Category  *CategoryRepository
	Document  *DocumentRepository
}

// New создаёт экземпляры всех репозиториев из переданного пула соединений.
func New(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:      NewUserRepository(pool),
		Regatta:   NewRegattaRepository(pool),
		Sailor:    NewSailorRepository(pool),
		Race:      NewRaceRepository(pool),
		RaceGroup: NewRaceGroupRepository(pool),
		Category:  NewCategoryRepository(pool),
		Document:  NewDocumentRepository(pool),
	}
}