package authentication

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication/mocks"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/stretchr/testify/require"
)

func TestValidToken(t *testing.T) {
	assert := require.New(t)
	repository := new(mocks.Repository)
	sso := New(repository)

	repository.On("GetUserByID", int64(123)).Return(&models.User{
		Id: 123,
	}, nil)

	userCtx, err := sso.ValidateToken("token-123")
	assert.NoError(err)
	assert.Equal(int64(123), userCtx.User.Id)
}

func TestInvalidToken(t *testing.T) {
	assert := require.New(t)
	repository := new(mocks.Repository)
	sso := New(repository)

	repository.On("GetUserByID", int64(123)).Return(nil, nil)

	userCtx, err := sso.ValidateToken("token-123")
	assert.Equal(err, ErrMissingUser)
	assert.Nil(userCtx)
}
