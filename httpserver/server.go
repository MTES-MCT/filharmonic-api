package httpserver

import (
	"log"
	"net/http"
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication"

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
	}

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

	tryCount := 100
	for tryCount > 0 {
		_, err := http.Get("http://" + s.config.Host + ":" + s.config.Port)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		tryCount--
	}
	if tryCount == 0 {
		log.Fatalln("server did not start")
	}

	return server
}
