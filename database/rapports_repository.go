package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) CreateRapport(idInspection int64, rapport models.Rapport) error {
	return repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		err := tx.Insert(&rapport)
		if err != nil {
			return err
		}
		inspection := models.Inspection{
			Id:        idInspection,
			RapportId: rapport.Id,
		}
		columns := []string{"rapport_id"}
		_, err = tx.Model(&inspection).Column(columns...).WherePK().Update()
		return err
	})
}

func (repo *Repository) GetRapport(ctx *domain.UserContext, idInspection int64) (*models.Rapport, error) {
	rapport := &models.Rapport{}
	query := repo.db.client.Model(rapport).
		Join("JOIN inspections AS inspection").
		JoinOn("inspection.rapport_id = rapport.id").
		JoinOn("inspection.id = ?", idInspection)

	if ctx.IsExploitant() {
		query.Join("JOIN etablissements AS etablissement").
			JoinOn("etablissement.id = inspection.etablissement_id").
			Join("JOIN etablissement_to_exploitants AS ex").
			JoinOn("ex.etablissement_id = etablissement.id").
			JoinOn("ex.user_id = ?", ctx.User.Id)
	} else {
		query.Join("JOIN inspection_to_inspecteurs AS u").
			JoinOn("u.inspection_id = inspection.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return rapport, err
}
