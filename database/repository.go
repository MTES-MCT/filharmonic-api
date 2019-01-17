package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

type Repository struct {
	db *Database
}

func NewRepository(db *Database) *Repository {
	return &Repository{
		db: db,
	}
}

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
