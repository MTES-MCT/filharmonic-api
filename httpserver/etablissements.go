package httpserver

import (
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) listEtablissements(c *gin.Context) {
	filter := domain.ListEtablissementsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	etablissements, err := server.service.ListEtablissements(server.retrieveUserContext(c), filter)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, etablissements)
}

func (server *HttpServer) getEtablissement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	etablissement, err := server.service.GetEtablissement(server.retrieveUserContext(c), id)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if etablissement == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, etablissement)
}
