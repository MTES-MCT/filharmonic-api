package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEscapeString(t *testing.T) {
	assert := require.New(t)
	assert.Equal("NomEtab", EscapeString("NomEtab"))
	assert.Equal("Nom\\%Et\\%\\%a\\\\%b", EscapeString("Nom%Et%%a\\%b"))
}
