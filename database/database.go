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
	Host       string `default:"localhost"`
	Port       int    `default:"5432"`
	User       string `default:"filharmonic"`
	Password   string `default:"filharmonic"`
	Name       string `default:"filharmonic"`
	InitSchema bool   `default:"false"`
	Seeds      bool   `default:"false"`
	LogSQL     bool   `default:"false"`
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
	&models.Inspection{},
	&models.Constat{},
	&models.PointDeControle{},
	&models.Theme{},
	&models.Commentaire{},
	&models.Message{},
	&models.PieceJointe{},
	&models.InspectionToInspecteur{},
	&models.UserToFavori{},
}

func New(config Config) (*Database, error) {
	client := pg.Connect(&pg.Options{
		User:     config.User,
		Password: config.Password,
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Database: config.Name,
	})
	if config.LogSQL {
		client.AddQueryHook(dbLogger{})
	}
	_, err := client.ExecOne("select 1")
	if err != nil {
		return nil, err
	}
	log.Info().Msg("connected to the database")

	db := &Database{
		config: config,
		client: client,
	}
	if config.InitSchema {
		err = db.createSchema()
	} else {
		err = migrations.MigrateDB(client)
	}
	if err != nil {
		return nil, err
	}
	if config.Seeds {
		err = seeds.SeedsTestDB(client)
	}
	return db, err
}

func (d *Database) createSchema() error {
	log.Warn().Msg("clearing and creating database")
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

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}
