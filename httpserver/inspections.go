package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) listInspections(c *gin.Context) {
	inspections, err := server.service.ListInspections(server.retrieveUserContext(c))
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, inspections)
}
