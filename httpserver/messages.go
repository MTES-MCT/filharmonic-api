package httpserver

import (
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) addMessage(c *gin.Context) {
	var message models.Message
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err = c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	idMessage, err := server.service.CreateMessage(server.retrieveUserContext(c), idPointDeControle, message)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": idMessage})
}

func (server *HttpServer) lireMessage(c *gin.Context) {
	idMessage, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = server.service.LireMessage(server.retrieveUserContext(c), idMessage)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "read"})
}
