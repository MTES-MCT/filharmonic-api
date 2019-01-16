package app

import (
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
)

type Config struct {
	Database database.Config
	Http     httpserver.Config
	LogLevel string `default:"info"`
}
