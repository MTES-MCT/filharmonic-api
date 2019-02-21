package httpserver

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listInspecteurs(c *gin.Context) (interface{}, error) {
	filters := domain.ListUsersFilters{
		Inspecteurs: true,
	}
	return server.service.ListUsers(server.retrieveUserContext(c), filters)
}
