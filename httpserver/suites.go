package httpserver

import (
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) addSuite(c *gin.Context) {
	var suite models.Suite
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"suite": err.Error(),
		})
		return
	}
	if err = c.ShouldBindJSON(&suite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	idSuite, err := server.service.CreateSuite(server.retrieveUserContext(c), idInspection, suite)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": idSuite})
}

func (server *HttpServer) updateSuite(c *gin.Context) {
	var suite models.Suite
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"suite": err.Error(),
		})
		return
	}
	if err = c.ShouldBindJSON(&suite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = server.service.UpdateSuite(server.retrieveUserContext(c), idInspection, suite)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (server *HttpServer) deleteSuite(c *gin.Context) {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = server.service.DeleteSuite(server.retrieveUserContext(c), idInspection)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
