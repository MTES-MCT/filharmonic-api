package authentication

import (
	"errors"
	"strconv"
	"strings"

	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
)

var (
	ErrUnauthorized = errors.New("unauthorized user")
	ErrMissingUser  = errors.New("missing user")
)

type Sso struct {
	repo Repository
}

func New(repo Repository) *Sso {
	return &Sso{
		repo: repo,
	}
}

func generateToken(id int64) string {
	return "token-" + strconv.FormatInt(id, 10)
}

func (sso *Sso) Login(email string, password string) (string, error) {
	user, err := sso.repo.GetUserByEmail(email)
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

func (sso *Sso) ValidateToken(token string) (int64, error) {
	userIdStr := strings.TrimLeft(token, "token-")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, err
	}
	user, err := sso.repo.GetUserByID(userId)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, ErrMissingUser
	}
	return user.ID, nil
}
