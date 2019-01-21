package domain

import "github.com/MTES-MCT/filharmonic-api/models"

//go:generate mockery -name Repository

type Repository interface {
	FindEtablissementsByS3IC(ctx *UserContext, s3ic string) ([]models.Etablissement, error)
	GetEtablissementByID(ctx *UserContext, id int64) (*models.Etablissement, error)
	ListInspections(ctx *UserContext) ([]models.Inspection, error)
	GetInspectionByID(ctx *UserContext, id int64) (*models.Inspection, error)

	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}
