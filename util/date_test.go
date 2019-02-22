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

func TestFormatDateTime(t *testing.T) {
	assert := require.New(t)
	assert.Equal("20/02/2019 à 08h40", FormatDateTime(DateTime("2019-02-20T08:40:06")))
	assert.Equal("20/12/2019 à 15h08", FormatDateTime(DateTime("2019-12-20T15:08:12")))
	assert.Equal("01/01/2001 à 23h59", FormatDateTime(DateTime("2001-01-01T23:59:59")))
}
