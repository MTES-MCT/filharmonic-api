package authentication

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/rs/zerolog"

	"github.com/MTES-MCT/filharmonic-api/authentication/mocks"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/stretchr/testify/require"
)

type AuthenticationServiceMocks struct {
	Repo     *mocks.Repository
	Sso      *mocks.Sso
	Sessions sessions.Sessions
}

func initAuthenticationService(t *testing.T) (*require.Assertions, *AuthenticationService, *AuthenticationServiceMocks) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	assert := require.New(t)
	repository := new(mocks.Repository)
	sso := new(mocks.Sso)
	sessions := sessions.NewMemory()
	mocks := AuthenticationServiceMocks{
		Repo:     repository,
		Sso:      sso,
		Sessions: sessions,
	}
	authenticationService := New(repository, sso, sessions)
	return assert, authenticationService, &mocks
}

func TestValidateTokenSuccessful(t *testing.T) {
	assert, authenticationService, mock := initAuthenticationService(t)

	mock.Sessions.Set("token-123", int64(123))
	mock.Repo.On("GetUserByID", int64(123)).Return(&models.User{
		Id: 123,
	}, nil)

	userCtx, err := authenticationService.ValidateToken("token-123")
	assert.NoError(err)
	assert.Equal(int64(123), userCtx.User.Id)
}

func TestValidateTokenFailed(t *testing.T) {
	assert, authenticationService, mock := initAuthenticationService(t)

	mock.Repo.On("GetUserByID", int64(123)).Return(nil, nil)

	userCtx, err := authenticationService.ValidateToken("token-123")
	assert.Equal(ErrUnauthorized, err)
	assert.Nil(userCtx)
}

func TestLoginSuccessful(t *testing.T) {
	assert, authenticationService, mock := initAuthenticationService(t)
	email := "inspecteur1@filharmonic.com"

	mock.Sessions.Set("token-123", int64(123))
	mock.Sso.On("ValidateTicket", "ticket-123").Return(email, nil)
	mock.Repo.On("GetUserByEmail", email).Return(&models.User{
		Id:    123,
		Email: email,
	}, nil)
	mock.Sessions.Set("token-123", int64(123))

	loginResult, err := authenticationService.Login("ticket-123")
	assert.NoError(err)
	assert.Equal(int64(123), loginResult.User.Id)
	assert.Equal(email, loginResult.User.Email)
}

func TestLoginFailedWithMissingUser(t *testing.T) {
	assert, authenticationService, mock := initAuthenticationService(t)
	email := "inspecteur1@filharmonic.com"

	mock.Sessions.Set("token-123", int64(123))
	mock.Sso.On("ValidateTicket", "ticket-123").Return(email, nil)
	mock.Repo.On("GetUserByEmail", email).Return(nil, nil)

	_, err := authenticationService.Login("ticket-123")
	assert.Error(err)
}
