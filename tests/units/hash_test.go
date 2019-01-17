package units

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/stretchr/testify/require"
)

func TestGeneratePassword(t *testing.T) {
	assert := require.New(t)
	hash, err := hash.GenerateFromPassword("toto")
	assert.NoError(err)
	assert.NotEmpty(hash)
}
func TestComparePasswordAndHash(t *testing.T) {
	assert := require.New(t)
	encodedhash, err := hash.GenerateFromPassword("toto")
	assert.NoError(err)
	checked, err := hash.ComparePasswordAndHash("toto", encodedhash)
	assert.NoError(err)
	assert.True(checked)
	checked, err = hash.ComparePasswordAndHash("tutu", encodedhash)
	assert.NoError(err)
	assert.False(checked)
}
