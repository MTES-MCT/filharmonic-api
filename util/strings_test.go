package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	assert := require.New(t)
	assert.Equal("equipement", Normalize("équipement"))
	assert.Equal("Equipement", Normalize("Équipement"))
	assert.Equal("hotel", Normalize("hôtel"))
	assert.Equal("A Nantes", Normalize("À Nantes"))
	assert.Equal("A-Nantes", Normalize("À-Nantes"))
	assert.Equal("A_Nantes", Normalize("À_Nantes"))
}
