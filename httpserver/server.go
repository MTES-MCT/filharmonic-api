package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/gin-contrib/logger"
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
	config                Config
	service               *domain.Service
	authenticationService *authentication.AuthenticationService
	server                *http.Server
}

func New(config Config, service *domain.Service, authenticationService *authentication.AuthenticationService) *HttpServer {
	return &HttpServer{
		config:                config,
		service:               service,
		authenticationService: authenticationService,
	}
}

func (s *HttpServer) Start() (returnErr error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	if s.config.Logger {
		router.Use(logger.SetLogger())
	}
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server up"})
	})
	router.POST("/login", s.login)
	router.POST("/authenticate", s.authenticate)

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
		authorized.POST("/inspections/:id/publier", s.publishInspection)
		authorized.POST("/inspections/:id/demandervalidation", s.askValidateInspection)
		authorized.POST("/inspections/:id/valider", s.validateInspection)
		authorized.POST("/inspections/:id/rejeter", s.rejectInspection)
		authorized.POST("/inspections/:id/suite", s.addSuite)
		authorized.PUT("/inspections/:id/suite", s.updateSuite)
		authorized.DELETE("/inspections/:id/suite", s.deleteSuite)
		authorized.PUT("/pointsdecontrole/:id", s.updatePointDeControle)
		authorized.POST("/pointsdecontrole/:id/constat", s.addConstat)
		authorized.DELETE("/pointsdecontrole/:id/constat", s.deleteConstat)
		authorized.POST("/pointsdecontrole/:id/publier", s.publishPointDeControle)
		authorized.POST("/pointsdecontrole/:id/messages", s.addMessage)
		authorized.DELETE("/pointsdecontrole/:id", s.deletePointDeControle)
		authorized.POST("/messages/:id/lire", s.lireMessage)
	}

	s.server = &http.Server{
		Addr:    s.config.Host + ":" + s.config.Port,
		Handler: router,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			returnErr = err
		}
	}()

	tryCount := 100
	for tryCount > 0 {
		_, err := http.Get("http://" + s.server.Addr)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		tryCount--
	}
	if tryCount == 0 {
		return errors.New("server did not start in time")
	}
	if returnErr == nil {
		log.Info().Msgf("server ready and listening on %s", s.server.Addr)
	}
	return returnErr
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
