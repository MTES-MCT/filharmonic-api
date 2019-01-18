package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

func MigrateDB(db *pg.DB) error {
	_, _, _ = migrations.Run(db, "init") // #nosec G104
	oldVersion, newVersion, err := migrations.Run(db, "up")
	if err != nil {
		return err
	}
	if newVersion != oldVersion {
		log.Info().Msgf("migrated db schema from version %d to %d", oldVersion, newVersion)
	} else {
		log.Info().Msgf("db schema version is up to date at %d", oldVersion)
	}
	return err
}
