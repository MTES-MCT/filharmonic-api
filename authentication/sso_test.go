package authentication

import (
	"encoding/xml"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication/mocks"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/stretchr/testify/require"
)

func initSSO(t *testing.T) (*require.Assertions, *Sso, *mocks.Repository) {
	assert := require.New(t)
	repository := new(mocks.Repository)
	ssoConfig := Config{
		// URL: ""
	}
	sso := New(ssoConfig, repository)
	return assert, sso, repository
}

func TestValidToken(t *testing.T) {
	assert, sso, repository := initSSO(t)

	repository.On("GetUserByID", int64(123)).Return(&models.User{
		Id: 123,
	}, nil)

	userCtx, err := sso.ValidateToken("token-123")
	assert.NoError(err)
	assert.Equal(int64(123), userCtx.User.Id)
}

func TestInvalidToken(t *testing.T) {
	assert, sso, repository := initSSO(t)

	repository.On("GetUserByID", int64(123)).Return(nil, nil)

	userCtx, err := sso.ValidateToken("token-123")
	assert.Equal(err, ErrMissingUser)
	assert.Nil(userCtx)
}

func TestUnmarshalValidateTicketResponse(t *testing.T) {
	assert := require.New(t)

	data := `<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas"><cas:authenticationSuccess><cas:user>user@example.com</cas:user></cas:authenticationSuccess></cas:serviceResponse>`

	userInfos := UserInfos{}
	err := xml.Unmarshal([]byte(data), &userInfos)
	assert.NoError(err)
	assert.Equal("user@example.com", userInfos.Email)
}
