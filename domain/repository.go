package domain

import "github.com/MTES-MCT/filharmonic-api/models"

//go:generate mockery -name Repository

type Repository interface {
	FindEtablissementsByS3IC(ctx *UserContext, s3ic string) ([]models.Etablissement, error)
	GetEtablissementByID(ctx *UserContext, id int64) (*models.Etablissement, error)

	ListInspections(ctx *UserContext) ([]models.Inspection, error)
	CreateInspection(ctx *UserContext, inspection models.Inspection) (int64, error)
	UpdateInspection(ctx *UserContext, inspection models.Inspection) error
	GetInspectionByID(ctx *UserContext, id int64) (*models.Inspection, error)

	CreatePointDeControle(ctx *UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error)
	UpdatePointDeControle(ctx *UserContext, idPointDeControle int64, pointDeControle models.PointDeControle) error
	DeletePointDeControle(ctx *UserContext, idPointDeControle int64) error
	PublishPointDeControle(ctx *UserContext, idPointDeControle int64) error

	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	CheckUsersInspecteurs(ids []int64) (bool, error)
	CheckInspecteurAllowedInspection(ctx *UserContext, idInspection int64) (bool, error)
	CheckInspecteurAllowedPointDeControle(ctx *UserContext, idPointDeControle int64) (bool, error)
}
