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
	router.POST("/login", returnResult(s.login))
	router.POST("/authenticate", returnResult(s.authenticate))

	authorized := router.Group("/")
	authorized.Use(s.authRequired)
	{
		authorized.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})
		authorized.POST("/logout", returnOk(s.logout))
		authorized.GET("/user", returnResult(s.getUser))
		authorized.GET("/inspecteurs", returnResult(s.listInspecteurs))
		authorized.GET("/etablissements", returnResult(s.listEtablissements))
		authorized.GET("/etablissements/:id", returnResult(s.getEtablissement))
		authorized.GET("/inspections", returnResult(s.listInspections))
		authorized.POST("/inspections", returnId(s.createInspection))
		authorized.GET("/inspectionsfavorites", returnResult(s.listInspectionsFavorites))
		authorized.GET("/inspections/:id", returnResult(s.getInspection))
		authorized.PUT("/inspections/:id", returnOk(s.updateInspection))
		authorized.POST("/inspections/:id/commentaires", returnId(s.addCommentaire))
		authorized.POST("/inspections/:id/pointsdecontrole", returnId(s.addPointDeControle))
		authorized.POST("/inspections/:id/publier", returnOk(s.publishInspection))
		authorized.POST("/inspections/:id/demandervalidation", returnOk(s.askValidateInspection))
		authorized.POST("/inspections/:id/valider", returnOk(s.validateInspection))
		authorized.POST("/inspections/:id/rejeter", returnOk(s.rejectInspection))
		authorized.POST("/inspections/:id/clore", returnOk(s.cloreInspection))
		authorized.POST("/inspections/:id/suite", returnId(s.addSuite))
		authorized.PUT("/inspections/:id/suite", returnOk(s.updateSuite))
		authorized.DELETE("/inspections/:id/suite", returnOk(s.deleteSuite))
		authorized.POST("/inspections/:id/favori", returnOk(s.addFavoriToInspection))
		authorized.DELETE("/inspections/:id/favori", returnOk(s.removeFavoriToInspection))
		authorized.GET("/inspections/:id/generer/lettreannonce", returnFile(s.genererLettreAnnonce))
		authorized.GET("/inspections/:id/generer/lettresuite", returnFile(s.genererLettreSuite))
		authorized.GET("/inspections/:id/generer/rapport", returnFile(s.genererRapport))
		authorized.GET("/inspections/:id/rapport", returnFile(s.getRapport))
		authorized.POST("/inspections/:id/canevas", returnId(s.createCanevas))
		authorized.PUT("/pointsdecontrole/:id", returnOk(s.updatePointDeControle))
		authorized.POST("/pointsdecontrole/:id/constat", returnId(s.addConstat))
		authorized.POST("/pointsdecontrole/:id/constat/resoudre", returnOk(s.resolveConstat))
		authorized.DELETE("/pointsdecontrole/:id/constat", returnOk(s.deleteConstat))
		authorized.POST("/pointsdecontrole/:id/publier", returnOk(s.publishPointDeControle))
		authorized.POST("/pointsdecontrole/:id/messages", returnId(s.addMessage))
		authorized.DELETE("/pointsdecontrole/:id", returnOk(s.deletePointDeControle))
		authorized.POST("/messages/:id/lire", returnOk(s.lireMessage))
		authorized.GET("/themes", returnResult(s.listThemes))
		authorized.POST("/themes", returnId(s.createTheme))
		authorized.DELETE("/themes/:id", returnOk(s.deleteTheme))
		authorized.GET("/canevas", returnResult(s.listCanevas))
		authorized.DELETE("/canevas/:id", returnOk(s.deleteCanevas))
		authorized.POST("/piecesjointes", returnId(s.createPieceJointe))
		authorized.GET("/piecesjointes/:id", returnFile(s.getPieceJointe))
		authorized.GET("/notifications", returnResult(s.listNotifications))
		authorized.POST("/notifications/lire", returnOk(s.lireNotifications))
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
