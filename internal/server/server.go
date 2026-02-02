package server

import "github.com/gin-gonic/gin"

type Server struct {
	router *gin.Engine
}

func New() *Server {
	return &Server{
		router: setupRouter(),
	}
}

// Run starts serving HTTP traffic.
func (s *Server) Run() error {
	return s.router.Run(":8080")
}

