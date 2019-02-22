package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatDate(t *testing.T) {
	assert := require.New(t)
	assert.Equal("20/02/2019", FormatDate(Date("2019-02-20").Time))
	assert.Equal("20/12/2019", FormatDate(Date("2019-12-20").Time))
	assert.Equal("01/01/2001", FormatDate(Date("2001-01-01").Time))
}
