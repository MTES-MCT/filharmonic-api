package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) addPointDeControle(c *gin.Context) (int64, error) {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorI(err)
	}
	var pointDeControle models.PointDeControle
	if err = c.ShouldBindJSON(&pointDeControle); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreatePointDeControle(server.retrieveUserContext(c), idInspection, pointDeControle)
}

func (server *HttpServer) updatePointDeControle(c *gin.Context) error {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	var pointDeControle models.PointDeControle
	if err = c.ShouldBindJSON(&pointDeControle); err != nil {
		return badInputError(err)
	}
	return server.service.UpdatePointDeControle(server.retrieveUserContext(c), idPointDeControle, pointDeControle)
}

func (server *HttpServer) deletePointDeControle(c *gin.Context) error {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.DeletePointDeControle(server.retrieveUserContext(c), idPointDeControle)
}

func (server *HttpServer) publishPointDeControle(c *gin.Context) error {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.PublishPointDeControle(server.retrieveUserContext(c), idPointDeControle)
}
