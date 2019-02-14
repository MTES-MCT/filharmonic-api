package database

type RepositoryConfig struct {
	PaginationSize int `default:"50"`
}

type Repository struct {
	config RepositoryConfig
	db     *Database
}

func NewRepository(config RepositoryConfig, db *Database) *Repository {
	return &Repository{
		config: config,
		db:     db,
	}
}
