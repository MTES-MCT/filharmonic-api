package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

const userContextKey = "userContextKey"
const AuthorizationHeader = "Authorization"

func (server *HttpServer) authRequired(c *gin.Context) {
	token := c.Request.Header.Get(AuthorizationHeader)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	}
	userContext, err := server.sso.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	}
	c.Set(userContextKey, userContext)
	c.Next()
}

func (server *HttpServer) retrieveUserContext(c *gin.Context) *domain.UserContext {
	ctx, _ := c.Get(userContextKey)
	return ctx.(*domain.UserContext)
}
