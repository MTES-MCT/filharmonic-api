package importcsv

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/rs/zerolog/log"
)

const (
	batchSize = 2000
)

func LoadEtablissementsCSV(filepath string, db *database.Database) (returnErr error) {
	indexesEtablissement := map[string]int{
		"codeBase":    0,
		"codeEtab":    1,
		"nom":         2,
		"region":      3,
		"departement": 4,
		"codeInsee":   5,
		"codePostal":  6,
		"activite":    8,
		"commune":     9,
		"seveso":      10,
		"regime":      11,
		"iedmtd":      13,
		"adresse1":    15,
		"adresse2":    16,
	}
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

	_, err = db.Exec("UPDATE pg_index SET indisready = false WHERE indrelid = (SELECT oid FROM pg_class WHERE relname = 'etablissements') AND indisunique = false")
	if err != nil {
		return err
	}
	defer func() {
		_, returnErr = db.Exec(`UPDATE pg_index SET indisready = true WHERE indrelid = (SELECT oid FROM pg_class WHERE relname = 'etablissements' AND indisunique = false);
											REINDEX TABLE etablissements`)

		log.Warn().Msgf("%d établissements importés", nbEtablissementsImportes)
	}()

	departements := []models.Departement{}
	err = db.Model(&departements).Column("id", "code_insee").Select()
	if err != nil {
		return err
	}

	codeInseeToDepartementId := make(map[string]int64)
	for _, departement := range departements {
		codeInseeToDepartementId[departement.CodeInsee] = departement.Id
	}

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
			regime := models.RegimeFromString(line[indexesEtablissement["regime"]])
			if regime == models.RegimeInconnu {
				log.Error().Msgf("régime `%s` inconnu pour l'établissement `%s`", line[indexesEtablissement["regime"]], line[indexesEtablissement["nom"]])
				continue
			}

			codeDepartement := util.GetCodeDepartementFromCodeInseeCommune(line[indexesEtablissement["codeInsee"]])
			departementId, ok := codeInseeToDepartementId[codeDepartement]
			if !ok {
				log.Warn().Msgf("département `%s` non trouvé pour l'établissement `%s`", codeDepartement, line[indexesEtablissement["nom"]])
			}

			etablissement := models.Etablissement{
				S3IC:          computeS3IC(line[indexesEtablissement["codeBase"]], line[indexesEtablissement["codeEtab"]]),
				Nom:           line[indexesEtablissement["nom"]],
				Raison:        line[indexesEtablissement["nom"]],
				Activite:      line[indexesEtablissement["activite"]],
				Seveso:        line[indexesEtablissement["seveso"]],
				Regime:        regime,
				Iedmtd:        toBool(line[indexesEtablissement["iedmtd"]]),
				Adresse1:      line[indexesEtablissement["adresse1"]],
				Adresse2:      line[indexesEtablissement["adresse2"]],
				CodePostal:    line[indexesEtablissement["codePostal"]],
				Commune:       line[indexesEtablissement["commune"]],
				DepartementId: departementId,
			}
			etablissements = append(etablissements, etablissement)
			iterations++
			nbEtablissementsImportes += 1
		}
		if len(etablissements) > 0 {
			_, err = db.Model(&etablissements).
				OnConflict("(s3ic) DO UPDATE").
				Insert()
			if err != nil {
				log.Error().Err(err).Msg("failed to save etablissements")
			}
			etablissements = etablissements[:0]
		}
	}
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
