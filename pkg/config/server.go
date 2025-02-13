package config

import "fmt"

type Server struct {
	Host string `envconfig:"SERVER_HOST" default:"localhost"`
	Port int    `envconfig:"SERVER_PORT" default:"8080"`
}

func (s *Server) Address() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}
