package authentication

import (
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
)

var (
	ErrUnauthorized = domain.NewErrUnauthorized("Accès non autorisé")
)

type SsoConfig struct {
	URL string `json:"url" default:"https://authentification.din.developpement-durable.gouv.fr/cas/public"`
}

type AuthenticationService struct {
	repo     Repository
	sso      Sso
	sessions sessions.Sessions
}

func New(repo Repository, sso Sso, sessions sessions.Sessions) *AuthenticationService {
	return &AuthenticationService{
		repo:     repo,
		sso:      sso,
		sessions: sessions,
	}
}

type LoginResult struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (service *AuthenticationService) Login(ticket string) (*LoginResult, error) {
	log.Debug().Str("ticket", ticket).Msg("login")

	email, err := service.sso.ValidateTicket(ticket)
	if err != nil {
		return nil, err
	}

	user, err := service.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	token, err := service.sessions.Add(user.Id)
	if err != nil {
		return nil, err
	}
	result := &LoginResult{
		Token: token,
		User:  *user,
	}
	return result, nil
}

func (service *AuthenticationService) ValidateToken(token string) (*domain.UserContext, error) {
	userId, err := service.sessions.Get(token)
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, ErrUnauthorized
	}
	user, err := service.repo.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return &domain.UserContext{User: user}, nil
}

func (service *AuthenticationService) Logout(token string) error {
	log.Debug().Str("token", token).Msg("logout")
	return service.sessions.Delete(token)
}
