package config

import (
	"fmt"
	"os"
	"strconv"
)

type Server struct {
	Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"SERVER_PORT" default:"8080"`
}

func (s *Server) Address() string {
	// Cloud Run sets PORT environment variable
	// Check PORT first, then fall back to SERVER_PORT
	port := s.Port
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	return fmt.Sprintf("%v:%v", s.Host, port)
}
