package icpe

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
)

const (
	indexCodeBase   = 0
	indexCodeEtab   = 1
	indexNom        = 2
	indexCodePostal = 6
	indexActivite   = 8
	indexCommune    = 9
	indexSeveso     = 10
	indexIedmtd     = 13
	indexAdresse1   = 15
	indexAdresse2   = 16

	batchSize = 2000
)

func LoadEtablissementsCSV(filepath string, db *database.Database) error {
	file, err := os.Open(filepath) // #nosec
	if err != nil {
		return err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	log.Warn().Msg("début de l'import")
	nbEtablissementsImportes := 0
	_, err = reader.Read()
	if err == io.EOF {
		return err
	}
	etablissements := make([]models.Etablissement, 0)
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
			etablissement := models.Etablissement{
				S3IC:     computeS3IC(line[indexCodeBase], line[indexCodeEtab]),
				Nom:      line[indexNom],
				Raison:   line[indexNom],
				Activite: line[indexActivite],
				Seveso:   line[indexSeveso],
				Iedmtd:   toBool(line[indexIedmtd]),
				Adresse:  line[indexAdresse1] + " " + line[indexAdresse2] + " " + line[indexCodePostal] + " " + line[indexCommune],
			}
			etablissements = append(etablissements, etablissement)
			iterations++
			nbEtablissementsImportes += 1
		}
		_, err = db.Model(&etablissements).
			OnConflict("(s3ic) DO UPDATE").
			Insert()
		if err != nil {
			log.Error().Err(err).Msg("failed to save etablissements")
		}
		etablissements = etablissements[:0]
	}
	log.Warn().Msgf("%d établissements importés", nbEtablissementsImportes)

	return nil
}

func computeS3IC(codebase string, codeetab string) string {
	s3ic := pad(codebase, 4) + "." + pad(codeetab, 5)
	return s3ic
}

func pad(value string, length int) string {
	padSize := length - len(value)
	for i := 0; i < padSize; i++ {
		value = "0" + value
	}
	return value
}

func toBool(value string) bool {
	return value == "Oui"
}
