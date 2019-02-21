package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) addConstat(c *gin.Context) (int64, error) {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorI(err)
	}
	var constat models.Constat
	if err = c.ShouldBindJSON(&constat); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateConstat(server.retrieveUserContext(c), idPointDeControle, constat)
}

func (server *HttpServer) deleteConstat(c *gin.Context) error {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.DeleteConstat(server.retrieveUserContext(c), idPointDeControle)
}
