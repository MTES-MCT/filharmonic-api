package domain

import "github.com/MTES-MCT/filharmonic-api/models"

//go:generate mockery -name Repository

type Repository interface {
	FindEtablissementsByS3IC(s3ic string) ([]models.Etablissement, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}
