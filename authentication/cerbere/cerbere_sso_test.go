package cerbere

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalValidateTicketResponse(t *testing.T) {
	assert := require.New(t)

	data := `<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas"><cas:authenticationSuccess><cas:user>user@example.com</cas:user></cas:authenticationSuccess></cas:serviceResponse>`

	userInfos := UserInfos{}
	err := xml.Unmarshal([]byte(data), &userInfos)
	assert.NoError(err)
	assert.Equal("user@example.com", userInfos.Email)
}
