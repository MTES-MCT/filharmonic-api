package authentication

import (
	"errors"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/database"
)

var (
	ErrUnauthorized = errors.New("unauthorized user")
	ErrMissingUser  = errors.New("missing user")
)

type Sso struct {
	repo *database.Repository
}

func New(repo *database.Repository) *Sso {
	return &Sso{
		repo: repo,
	}
}

func generateToken(id int64) string {
	return "token-" + strconv.FormatInt(id, 10)
}

func (sso *Sso) Login(email string, password string) (string, error) {
	user, err := sso.repo.GetUser(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrMissingUser
	}
	checked, err := hash.ComparePasswordAndHash(password, user.Password)
	if err != nil {
		return "", err
	}
	if checked {
		return generateToken(user.ID), nil
	}
	return "", ErrUnauthorized
}
