package httpserver

import (
	"log"
	"net/http"
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string `default:"localhost"`
	Port string `default:"5000"`
}

type HttpServer struct {
	config  Config
	service *domain.Service
}

func New(config Config, service *domain.Service) *HttpServer {
	return &HttpServer{
		config:  config,
		service: service,
	}
}

func (s *HttpServer) Start() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/etablissements", s.listEtablissements)

	server := &http.Server{
		Addr:    s.config.Host + ":" + s.config.Port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	time.Sleep(time.Second)
	return server
}
