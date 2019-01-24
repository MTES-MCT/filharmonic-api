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
	for {
		line, err := reader.Read()
		if err == io.EOF {
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

		_, err = db.Model(&etablissement).
			OnConflict("(s3ic) DO UPDATE").
			Insert()
		if err != nil {
			log.Error().Interface("etablissement", etablissement).Msg("failed to save etablissement")
		} else {
			nbEtablissementsImportes += 1
		}
	}
	log.Warn().Msgf("%d établissements importés\n", nbEtablissementsImportes)

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
