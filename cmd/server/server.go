package main

import "github.com/gin-gonic/gin"

type Server struct {
	r *gin.Engine
	m *TCPManager
}

func InitServer() *Server {
	var s Server
	s.r = InitHttpRouter()
	s.m = &TCPManager{}

	return &s
}

func (s *Server) Run() {
	//err := s.r.Run(":7890")
	s.m.Run()
}
