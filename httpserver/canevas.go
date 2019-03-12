package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listCanevas(c *gin.Context) (interface{}, error) {
	return server.service.ListCanevas(server.retrieveUserContext(c))
}

func (server *HttpServer) createCanevas(c *gin.Context) (int64, error) {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorI(err)
	}
	var canevas models.Canevas
	if err := c.ShouldBindJSON(&canevas); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateCanevas(server.retrieveUserContext(c), idInspection, canevas)
}

func (server *HttpServer) deleteCanevas(c *gin.Context) error {
	idCanevas, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.DeleteCanevas(server.retrieveUserContext(c), idCanevas)
}
