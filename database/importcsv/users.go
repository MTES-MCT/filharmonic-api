package importcsv

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
)

func LoadInspecteursCSV(filepath string, db *database.Database) error {
	indexesInspecteur := map[string]int{
		"nom":      0,
		"prenom":   1,
		"nomUsage": 2,
		"email":    3,
	}
	file, err := os.Open(filepath) // #nosec
	if err != nil {
		return err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	log.Warn().Msg("début de l'import")
	nbInspecteursImportes := 0
	_, err = reader.Read()
	if err == io.EOF {
		return err
	}
	inspecteurs := make([]models.User, 0)
	done := false
	for !done {
		iterations := 0
		for iterations < batchSize {
			var line []string
			line, err = reader.Read()
			if err == io.EOF {
				done = true
				break
			}
			if err != nil {
				return err
			}
			nom := line[indexesInspecteur["nom"]]
			if line[indexesInspecteur["nomUsage"]] != "" {
				nom = line[indexesInspecteur["nomUsage"]]
			}
			inspecteur := models.User{
				Nom:     nom,
				Prenom:  line[indexesInspecteur["prenom"]],
				Email:   strings.ToLower(line[indexesInspecteur["email"]]),
				Profile: models.ProfilInspecteur,
			}
			inspecteurs = append(inspecteurs, inspecteur)
			iterations++
			nbInspecteursImportes += 1
		}
		_, err = db.Model(&inspecteurs).
			OnConflict("(email) DO UPDATE").
			Insert()
		if err != nil {
			log.Error().Err(err).Msg("failed to save inspecteurs")
		}
		inspecteurs = inspecteurs[:0]
	}
	log.Warn().Msgf("%d inspecteurs importés", nbInspecteursImportes)

	return nil
}
