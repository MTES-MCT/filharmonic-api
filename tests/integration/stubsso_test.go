package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/stubsso"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestValidateTicketSuccessful(t *testing.T) {
	assert, application := tests.InitDB(t)
	sso := stubsso.New(application.Repo)

	email, err := sso.ValidateTicket("ticket-1")
	assert.NoError(err)
	assert.Equal("exploitant1@filharmonic.com", email)
}

func TestValidateTicketFailed(t *testing.T) {
	assert, application := tests.InitDB(t)
	sso := stubsso.New(application.Repo)

	email, err := sso.ValidateTicket("ticket-9000000000000000000")
	assert.Equal(authentication.ErrMissingUser, err)
	assert.Equal("", email)
}
