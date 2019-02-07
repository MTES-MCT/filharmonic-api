package importcsv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPad(t *testing.T) {
	assert := require.New(t)
	assert.Equal("0023", pad("23", 4))
	assert.Equal("00023", pad("23", 5))
	assert.Equal("23", pad("23", 2))
	assert.Equal("23", pad("23", 1))
}
