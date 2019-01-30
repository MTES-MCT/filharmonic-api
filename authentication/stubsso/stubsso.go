package stubsso

import (
	"strconv"
	"strings"

	"github.com/MTES-MCT/filharmonic-api/authentication"

	"github.com/MTES-MCT/filharmonic-api/database"
)

type StubSso struct {
	repo *database.Repository
}

func New(repo *database.Repository) *StubSso {
	return &StubSso{
		repo: repo,
	}
}

func (sso *StubSso) ValidateTicket(ticket string) (string, error) {
	userIdStr := strings.TrimLeft(ticket, "ticket-")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return "", err
	}
	user, err := sso.repo.GetUserByID(userId)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", authentication.ErrMissingUser
	}
	return user.Email, nil
}
