package httpserver

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listEvenements(c *gin.Context) (interface{}, error) {
	filter := domain.ListEvenementsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		return badInputErrorN(err)
	}
	return server.service.ListEvenements(server.retrieveUserContext(c), filter)
}
