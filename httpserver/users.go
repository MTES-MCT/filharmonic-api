package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) listInspecteurs(c *gin.Context) {
	filters := domain.ListUsersFilters{
		Inspecteurs: true,
	}
	users, err := server.service.ListUsers(server.retrieveUserContext(c), filters)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}
