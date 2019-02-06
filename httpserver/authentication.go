package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const userContextKey = "userContextKey"
const AuthorizationHeader = "Authorization"

func (server *HttpServer) authRequired(c *gin.Context) {
	token := c.Request.Header.Get(AuthorizationHeader)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		log.Warn().Msg("token empty")
		return
	}
	userContext, err := server.authenticationService.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		log.Error().Err(err)
		return
	}
	c.Set(userContextKey, userContext)
}

func (server *HttpServer) retrieveUserContext(c *gin.Context) *domain.UserContext {
	ctx, _ := c.Get(userContextKey)
	return ctx.(*domain.UserContext)
}
