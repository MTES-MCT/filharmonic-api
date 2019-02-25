package domain

import (
	"io"

	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/MTES-MCT/filharmonic-api/models"
)

//go:generate mockery -name Repository

type Repository interface {
	FindEtablissements(ctx *UserContext, filter ListEtablissementsFilter) (*models.FindEtablissementResults, error)
	GetEtablissementByID(ctx *UserContext, id int64) (*models.Etablissement, error)

	ListThemes() ([]models.Theme, error)
	CreateTheme(theme models.Theme) (int64, error)
	DeleteTheme(idTheme int64) error

	ListInspections(ctx *UserContext, filter ListInspectionsFilter) ([]models.Inspection, error)
	ListInspectionsFavorites(ctx *UserContext) ([]models.Inspection, error)
	CreateInspection(ctx *UserContext, inspection models.Inspection) (int64, error)
	UpdateInspection(ctx *UserContext, inspection models.Inspection) error
	GetInspectionByID(ctx *UserContext, id int64, filter InspectionFilter) (*models.Inspection, error)
	CheckEtatInspection(idInspection int64, etats []models.EtatInspection) (bool, error)
	UpdateEtatInspection(ctx *UserContext, id int64, etat models.EtatInspection) error
	AddFavoriToInspection(ctx *UserContext, idInspection int64) error
	RemoveFavoriToInspection(ctx *UserContext, idInspection int64) error
	CheckCanCreateSuite(ctx *UserContext, idInspection int64) (bool, error)
	CheckCanDeleteSuite(ctx *UserContext, idInspection int64) (bool, error)

	CreateSuite(ctx *UserContext, idInspection int64, suite models.Suite) (int64, error)
	UpdateSuite(ctx *UserContext, idInspection int64, suite models.Suite) error
	DeleteSuite(ctx *UserContext, idInspection int64) error

	CreateCommentaire(ctx *UserContext, idInspection int64, commentaire models.Commentaire) (int64, error)

	CreatePointDeControle(ctx *UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error)
	UpdatePointDeControle(ctx *UserContext, idPointDeControle int64, pointDeControle models.PointDeControle) error
	DeletePointDeControle(ctx *UserContext, idPointDeControle int64) error
	PublishPointDeControle(ctx *UserContext, idPointDeControle int64) error
	CreateMessage(ctx *UserContext, idPointDeControle int64, message models.Message) (int64, error)
	LireMessage(ctx *UserContext, idMessage int64) error
	CheckUserAllowedMessage(ctx *UserContext, idMessage int64) (bool, error)
	CheckUserIsRecipient(ctx *UserContext, idMessage int64) (bool, error)
	CheckCanCreateConstat(ctx *UserContext, idPointDeControle int64) (bool, error)
	CheckCanDeleteConstat(ctx *UserContext, idPointDeControle int64) (bool, error)

	CreateConstat(ctx *UserContext, idPointDeControle int64, constat models.Constat) (int64, error)
	DeleteConstat(ctx *UserContext, idPointDeControle int64) error

	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	CheckUsersInspecteurs(ids []int64) (bool, error)
	CheckInspecteurAllowedInspection(ctx *UserContext, idInspection int64) (bool, error)
	CheckUserAllowedPointDeControle(ctx *UserContext, idPointDeControle int64) (bool, error)
	FindUsers(filters ListUsersFilters) ([]models.User, error)

	GetPieceJointe(ctx *UserContext, id int64) (*models.PieceJointe, error)
	CreatePieceJointe(pieceJointe models.PieceJointe) (int64, error)

	ListNotifications(ctx *UserContext, filter *ListNotificationsFilter) ([]models.Notification, error)
	LireNotifications(ctx *UserContext, ids []int64) error
	ListEvenements(ctx *UserContext, filter ListEvenementsFilter) ([]models.Evenement, error)
	CreateEvenement(ctx *UserContext, typeEvenement models.TypeEvenement, idInspection int64, data map[string]interface{}) error
	ListNouveauxMessages() ([]NouveauxMessagesUser, error)
}

type Storage interface {
	Get(id string) (io.Reader, error)
	Put(file models.PieceJointeFile) (string, error)
}

type EmailService interface {
	Send(email emails.Email) error
}

type TemplateService interface {
	RenderHTMLEmailNouveauxMessages(data interface{}) (string, error)
	RenderLettreAnnonce(data interface{}) (string, error)
	RenderLettreSuite(data interface{}) (string, error)
	RenderRapport(data interface{}) (string, error)
}
