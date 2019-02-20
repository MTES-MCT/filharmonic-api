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
		"nom":         0,
		"prenom":      1,
		"nomUsage":    2,
		"email":       3,
		"approbateur": 5,
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
			profil := models.ProfilInspecteur
			if line[indexesInspecteur["approbateur"]] == "oui" {
				profil = models.ProfilApprobateur
			}
			inspecteur := models.User{
				Nom:     nom,
				Prenom:  line[indexesInspecteur["prenom"]],
				Email:   strings.ToLower(line[indexesInspecteur["email"]]),
				Profile: profil,
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

// nolint: gocyclo
func LoadExploitantsCSV(filepath string, db *database.Database) error {
	indexesExploitant := map[string]int{
		"nom":    0,
		"prenom": 1,
		"email":  2,
		"s3ic":   3,
	}
	file, err := os.Open(filepath) // #nosec
	if err != nil {
		return err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	log.Warn().Msg("début de l'import")
	nbExploitantsImportes := 0
	_, err = reader.Read()
	if err == io.EOF {
		return err
	}
	etablissements := []models.Etablissement{}
	err = db.Model(&etablissements).Column("id", "s3ic").Select()
	if err != nil {
		return err
	}

	s3icToEtablissementId := make(map[string]int64)
	for _, etablissement := range etablissements {
		s3icToEtablissementId[etablissement.S3IC] = etablissement.Id
	}

	_, err = db.Model(&models.EtablissementToExploitant{}).Where("1 = 1").Delete()
	if err != nil {
		return err
	}

	exploitants := make([]models.User, 0)
	etablissementToExploitants := make([]models.EtablissementToExploitant, 0)
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
			s3ic := line[indexesExploitant["s3ic"]]
			etablissementId, ok := s3icToEtablissementId[s3ic]
			if !ok {
				log.Warn().Str("s3ic", s3ic).Msg("s3ic non trouvé")
				continue
			}

			exploitant := models.User{
				Nom:     line[indexesExploitant["nom"]],
				Prenom:  line[indexesExploitant["prenom"]],
				Email:   strings.ToLower(line[indexesExploitant["email"]]),
				Profile: models.ProfilExploitant,
			}
			exploitants = append(exploitants, exploitant)
			etablissementToExploitants = append(etablissementToExploitants, models.EtablissementToExploitant{
				EtablissementId: etablissementId,
			})

			iterations++
			nbExploitantsImportes += 1
		}
		_, err = db.Model(&exploitants).
			OnConflict("(email) DO UPDATE").
			Insert()
		if err != nil {
			log.Error().Err(err).Msg("failed to save exploitants")
		}

		for index, _ := range etablissementToExploitants {
			etablissementToExploitants[index].UserId = exploitants[index].Id
		}

		_, err = db.Model(&etablissementToExploitants).Insert()
		if err != nil {
			log.Error().Err(err).Msg("failed to save etablissementToExploitants")
		}

		exploitants = exploitants[:0]
		etablissementToExploitants = etablissementToExploitants[:0]
	}
	log.Warn().Msgf("%d exploitants importés", nbExploitantsImportes)

	return nil
}
