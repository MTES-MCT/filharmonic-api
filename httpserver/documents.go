package httpserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) genererLettreAnnonce(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	lettreAnnonceFile, err := server.service.GenererLettreAnnonce(server.retrieveUserContext(c), id)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	headers := make(map[string]string)
	c.Render(http.StatusOK, render.Reader{
		Headers:       headers,
		ContentType:   lettreAnnonceFile.Type,
		ContentLength: lettreAnnonceFile.Taille,
		Reader:        lettreAnnonceFile.Content,
	})
}
