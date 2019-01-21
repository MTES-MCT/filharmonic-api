package database

import (
	"fmt"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/database/migrations"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Config struct {
	Host       string `default:"localhost"`
	Port       int    `default:"5432"`
	User       string `default:"filharmonic"`
	Password   string `default:"filharmonic"`
	Name       string `default:"filharmonic"`
	InitSchema bool   `default:"false"`
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
	&models.Inspection{},
	&models.ThemeInspection{},
	&models.ThemeReferentiel{},
	&models.InspectionToInspecteur{},
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
	db := &Database{
		config: config,
		client: client,
	}
	if config.InitSchema {
		err = db.createSchema()
	} else {
		err = migrations.MigrateDB(client)
	}
	return db, err
}

func (d *Database) createSchema() error {
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

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}
