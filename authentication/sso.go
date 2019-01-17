package authentication

import (
	"errors"
	"strconv"
	"strings"

	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/domain"
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
		return generateToken(user.Id), nil
	}
	return "", ErrUnauthorized
}

func (sso *Sso) ValidateToken(token string) (*domain.UserContext, error) {
	userIdStr := strings.TrimLeft(token, "token-")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := sso.repo.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrMissingUser
	}
	return &domain.UserContext{User: user}, nil
}
