package httpserver

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listNotifications(c *gin.Context) (interface{}, error) {
	filter := domain.ListNotificationsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		return badInputErrorN(err)
	}
	return server.service.ListNotifications(server.retrieveUserContext(c), &filter)
}

func (server *HttpServer) lireNotifications(c *gin.Context) error {
	ids := []int64{}
	if err := c.ShouldBindJSON(&ids); err != nil {
		return badInputError(err)
	}
	return server.service.LireNotifications(server.retrieveUserContext(c), ids)
}
