package domain

import "github.com/MTES-MCT/filharmonic-api/models"

type Repository interface {
	FindEtablissementsByS3IC(s3ic string) ([]models.Etablissement, error)
	GetUser(email string) (*models.User, error)
}
