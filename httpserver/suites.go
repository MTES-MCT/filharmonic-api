package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) addSuite(c *gin.Context) (int64, error) {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorI(err)
	}
	var suite models.Suite
	if err = c.ShouldBindJSON(&suite); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateSuite(server.retrieveUserContext(c), idInspection, suite)
}

func (server *HttpServer) updateSuite(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	var suite models.Suite
	if err = c.ShouldBindJSON(&suite); err != nil {
		return badInputError(err)
	}
	return server.service.UpdateSuite(server.retrieveUserContext(c), idInspection, suite)
}

func (server *HttpServer) deleteSuite(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.DeleteSuite(server.retrieveUserContext(c), idInspection)
}
