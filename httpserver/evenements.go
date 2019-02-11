package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) listEvenements(c *gin.Context) {
	filter := domain.ListEvenementsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	evenements, err := server.service.ListEvenements(server.retrieveUserContext(c), filter)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, evenements)
}
