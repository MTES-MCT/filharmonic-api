package httpserver

import (
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) addConstat(c *gin.Context) {
	var constat models.Constat
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"constat": err.Error(),
		})
		return
	}
	if err = c.ShouldBindJSON(&constat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	idConstat, err := server.service.CreateConstat(server.retrieveUserContext(c), idPointDeControle, constat)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": idConstat})
}

func (server *HttpServer) deleteConstat(c *gin.Context) {
	idPointDeControle, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = server.service.DeleteConstat(server.retrieveUserContext(c), idPointDeControle)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
