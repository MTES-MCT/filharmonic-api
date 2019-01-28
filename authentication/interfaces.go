package authentication

import "github.com/MTES-MCT/filharmonic-api/models"

//go:generate mockery -all

type Repository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}

type Sso interface {
	ValidateTicket(ticket string) (string, error)
}
