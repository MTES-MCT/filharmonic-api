package database

import (
	"fmt"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/database/migrations"
	"github.com/MTES-MCT/filharmonic-api/database/seeds"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host            string `default:"localhost"`
	Port            int    `default:"5432"`
	User            string `default:"filharmonic"`
	Password        string `default:"filharmonic"`
	Name            string `default:"filharmonic"`
	ApplyMigrations bool   `default:"true"`
	InitSchema      bool   `default:"false"`
	Seeds           bool   `default:"false"`
	LogSQL          bool   `default:"false"`
	PaginationSize  int    `default:"50"`
}

type Database struct {
	config Config
	client *pg.DB
}

var tables = []interface{}{
	&models.Etablissement{},
	&models.User{},
	&models.EtablissementToExploitant{},
	&models.Suite{},
	&models.Rapport{},
	&models.Inspection{},
	&models.Constat{},
	&models.PointDeControle{},
	&models.Theme{},
	&models.Commentaire{},
	&models.Message{},
	&models.PieceJointe{},
	&models.InspectionToInspecteur{},
	&models.UserToFavori{},
	&models.Evenement{},
	&models.Notification{},
}

const createIndexesQuery = `CREATE EXTENSION pg_trgm;
CREATE INDEX trgm_idx_etablissments_s3ic ON etablissements USING gin (s3ic gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_nom ON etablissements USING gin (nom gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_raison ON etablissements USING gin (raison gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_adresse1 ON etablissements USING gin (adresse1 gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_adresse2 ON etablissements USING gin (adresse2 gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_code_postal ON etablissements USING gin (code_postal gin_trgm_ops);
CREATE INDEX trgm_idx_etablissments_commune ON etablissements USING gin (commune gin_trgm_ops);

CREATE INDEX idx_users_profile ON users(profile);
CREATE INDEX idx_inspections_etat ON inspections(etat);
CREATE INDEX idx_point_de_controles_publie ON point_de_controles(publie);
CREATE INDEX idx_messages_interne ON messages(interne);
CREATE INDEX idx_messages_lu ON messages(lu);
`

func New(config Config) (*Database, error) {
	address := config.Host + ":" + strconv.Itoa(config.Port)
	client := pg.Connect(&pg.Options{
		User:     config.User,
		Password: config.Password,
		Addr:     address,
		Database: config.Name,
	})
	if config.LogSQL {
		client.AddQueryHook(dbLogger{})
	}
	log.Info().Msgf("connecting to postgresql endpoint on %s", address)
	_, err := client.ExecOne("select 1")
	if err != nil {
		return nil, err
	}
	log.Info().Msg("connected to postgresql")

	db := &Database{
		config: config,
		client: client,
	}
	if config.InitSchema {
		log.Warn().Msg("clearing and creating database schema")
		err = db.createSchema()
		if err != nil {
			return nil, err
		}
		// Création des indexes cf migration 21_add_indexes_etablissements
		_, err = db.client.Exec(createIndexesQuery)
	} else if config.ApplyMigrations {
		err = migrations.MigrateDB(client)
	}
	if err != nil {
		return nil, err
	}
	if config.Seeds {
		log.Warn().Msg("seeding database")
		err = seeds.SeedsTestDB(client)
	}
	return db, err
}

func (d *Database) createSchema() error {
	_, err := d.client.Exec("drop extension if exists pg_trgm cascade")
	if err != nil {
		return err
	}
	for _, table := range tables {
		err := d.client.DropTable(table, &orm.DropTableOptions{
			Cascade:  true,
			IfExists: true,
		})
		if err != nil {
			return err
		}
	}
	for _, table := range tables {
		err := d.client.CreateTable(table, &orm.CreateTableOptions{
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) Close() error {
	return d.client.Close()
}

func (d *Database) Insert(model ...interface{}) error {
	return d.client.Insert(model...)
}

func (d *Database) Model(model ...interface{}) *orm.Query {
	return d.client.Model(model...)
}

func (d *Database) Exec(query interface{}, params ...interface{}) (pg.Result, error) {
	return d.client.Exec(query, params...)
}

func (d *Database) Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	return d.client.Query(model, query, params...)
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}
