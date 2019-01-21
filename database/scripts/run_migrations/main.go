package main

import (
	"log"

	"github.com/MTES-MCT/filharmonic-api/database"
)

func main() {
	db, err := database.New(database.Config{
		Host:       "localhost",
		Port:       5432,
		Name:       "filharmonic",
		User:       "filharmonic",
		Password:   "filharmonic",
		LogSQL:     true,
		InitSchema: false,
	})
	if err != nil {
		log.Fatalf("unable to connect to db, error: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal("Unable to close DB", err)
		}
	}()
}
