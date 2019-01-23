package httpserver

import (
	"net/http"
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/rs/zerolog/log"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Host   string `default:"0.0.0.0"`
	Port   string `default:"5000"`
	Logger bool   `default:"true"`
}

type HttpServer struct {
	config  Config
	service *domain.Service
	sso     *authentication.Sso
}

func New(config Config, service *domain.Service, sso *authentication.Sso) *HttpServer {
	return &HttpServer{
		config:  config,
		service: service,
		sso:     sso,
	}
}

func (s *HttpServer) Start() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	if s.config.Logger {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server up"})
	})
	router.POST("/login", s.login)

	authorized := router.Group("/")
	authorized.Use(s.authRequired)
	{
		authorized.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})
		authorized.GET("/etablissements", s.listEtablissements)
		authorized.GET("/etablissements/:id", s.getEtablissement)
		authorized.GET("/inspections", s.listInspections)
		authorized.POST("/inspections", s.createInspection)
		authorized.GET("/inspections/:id", s.getInspection)
		authorized.PUT("/inspections/:id", s.updateInspection)
		authorized.POST("/inspections/:id/commentaires", s.addCommentaire)
		authorized.POST("/inspections/:id/pointsdecontrole", s.addPointDeControle)
		authorized.PUT("/pointsdecontrole/:id", s.updatePointDeControle)
		authorized.POST("/pointsdecontrole/:id/publier", s.publishPointDeControle)
		authorized.POST("/pointsdecontrole/:id/messages", s.addMessage)
		authorized.DELETE("/pointsdecontrole/:id", s.deletePointDeControle)
		authorized.POST("/messages/:id/lire", s.lireMessage)
	}

	server := &http.Server{
		Addr:    s.config.Host + ":" + s.config.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen error: %s", err)
		}
	}()

	tryCount := 100
	for tryCount > 0 {
		_, err := http.Get("http://" + server.Addr)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		tryCount--
	}
	if tryCount == 0 {
		log.Fatal().Msg("server did not start")
	}
	log.Info().Msgf("server ready and listening on %s", server.Addr)
	return server
}
