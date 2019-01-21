package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) FindEtablissementsByS3IC(ctx *domain.UserContext, s3ic string) ([]models.Etablissement, error) {
	var etablissements []models.Etablissement
	query := repo.db.client.Model(&etablissements).Where("s3ic like ?", "%"+s3ic+"%")
	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	return etablissements, err
}

func (repo *Repository) GetEtablissementByID(ctx *domain.UserContext, id int64) (*models.Etablissement, error) {
	var etablissement models.Etablissement
	query := repo.db.client.Model(&etablissement).Where("id = ?", id)
	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &etablissement, err
}
