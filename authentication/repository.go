package authentication

import "github.com/MTES-MCT/filharmonic-api/models"

//go:generate mockery -name Repository

type Repository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}
