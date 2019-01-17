package tests

import (
	"testing"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/stretchr/testify/require"
)

func Init(t *testing.T, initDbFunc func(db *database.Database, assert *require.Assertions)) (*httpexpect.Expect, func()) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.Database.InitSchema = true
	config.Http.Host = "localhost"
	db, server := app.Bootstrap(config)
	if initDbFunc != nil {
		initDbFunc(db, assert)
	}
	e := httpexpect.New(t, "http://"+config.Http.Host+":"+config.Http.Port+"/")
	return e, func() {
		server.Close()
		// db.Close()
	}
}
