package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) addMessage(c *gin.Context) (int64, error) {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorI(err)
	}
	var message models.Message
	if err = c.ShouldBindJSON(&message); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateMessage(server.retrieveUserContext(c), idPointDeControle, message)
}

func (server *HttpServer) lireMessage(c *gin.Context) error {
	idMessage, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.LireMessage(server.retrieveUserContext(c), idMessage)
}
