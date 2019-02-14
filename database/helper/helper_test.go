package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEscapeString(t *testing.T) {
	assert := require.New(t)
	assert.Equal("%NomEtab%", BuildSearchValue("NomEtab"))
	assert.Equal("%Nom\\%Et\\%\\%a\\\\%b%", BuildSearchValue("Nom%Et%%a\\%b"))
}
