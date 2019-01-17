package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) listEtablissements(c *gin.Context) {
	s3ic := c.Query("s3ic")
	etablissements, err := server.service.ListEtablissements(server.retrieveUserContext(c), s3ic)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, etablissements)
}
