package database

import "github.com/MTES-MCT/filharmonic-api/models"

type Repository struct {
	db *Database
}

func NewRepository(db *Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) FindEtablissementsByS3IC(s3ic string) ([]models.Etablissement, error) {
	var etablissements []models.Etablissement
	err := repo.db.client.Model(&etablissements).Where("s3ic like ?", "%"+s3ic+"%").Select()
	return etablissements, err
}
