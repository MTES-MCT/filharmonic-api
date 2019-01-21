package database

type Repository struct {
	db *Database
}

func NewRepository(db *Database) *Repository {
	return &Repository{
		db: db,
	}
}
