package server

import (
	"github.com/dollarkillerx/analysis_japanese_schools/internal/conf"
	"github.com/gin-gonic/gin"
)

type Server struct {
	app *gin.Engine
}

func NewServer() *Server {
	eng := gin.New()
	eng.Use(gin.Recovery())
	eng.Use(gin.Logger())

	return &Server{
		app: eng,
	}
}

func (s *Server) Run() error {
	return s.app.Run(conf.GetConfig().ListenAddr)
}
