package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Credentials struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (server *HttpServer) login(c *gin.Context) {
	var user Credentials
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := server.sso.Login(user.Email, user.Password)
	if err != nil {
		if err != authentication.ErrMissingUser {
			log.Error().Err(err).Msg("Bad service response")
		}
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
