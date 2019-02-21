package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listEtablissements(c *gin.Context) (interface{}, error) {
	filter := domain.ListEtablissementsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		return badInputErrorN(err)
	}
	return server.service.ListEtablissements(server.retrieveUserContext(c), filter)
}

func (server *HttpServer) getEtablissement(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorN(err)
	}
	return server.service.GetEtablissement(server.retrieveUserContext(c), id)
}
