package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePassword(t *testing.T) {
	assert := require.New(t)
	hash, err := GenerateFromPassword("toto")
	assert.NoError(err)
	assert.NotEmpty(hash)
}
func TestComparePasswordAndHash(t *testing.T) {
	assert := require.New(t)
	encodedhash, err := GenerateFromPassword("toto")
	assert.NoError(err)
	checked, err := ComparePasswordAndHash("toto", encodedhash)
	assert.NoError(err)
	assert.True(checked)
	checked, err = ComparePasswordAndHash("tutu", encodedhash)
	assert.NoError(err)
	assert.False(checked)
}
