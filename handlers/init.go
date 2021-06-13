package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	cors "github.com/itsjamie/gin-cors"
)

type Options struct {
	Readonly bool
}

type HTTPServer struct {
	Server *http.Server
	Router *gin.Engine
	Opt    Options
	Repo   map[string]string
}

func NewHTTPServer(opt Options, repo map[string]string) (*HTTPServer, error) {
	logrus.Infof("Initializing gin server")

	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST",
		RequestHeaders:  "Origin, Content-Type, Authorization",
		ExposedHeaders:  "",
		MaxAge:          24 * 3600 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	h := &HTTPServer{
		Server: &http.Server{
			Addr:    ":3000",
			Handler: router,
		},
		Router: router,
		Opt:    opt,
		Repo:   repo,
	}

	h.setupRepoHandlers()

	return h, nil
}

//Start the main HTTP Server entry
func (s *HTTPServer) Start() error {
	logrus.Infof("Starting HTTP Server on port 3000")
	return s.Server.ListenAndServe()
}
