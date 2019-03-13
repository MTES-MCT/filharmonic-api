package database

import "github.com/MTES-MCT/filharmonic-api/events"

type RepositoryConfig struct {
	PaginationSize int `default:"50"`
}

type Repository struct {
	config        RepositoryConfig
	db            *Database
	eventsManager *events.EventsManager
}

func NewRepository(config RepositoryConfig, db *Database, eventsManager *events.EventsManager) *Repository {
	return &Repository{
		config:        config,
		db:            db,
		eventsManager: eventsManager,
	}
}
