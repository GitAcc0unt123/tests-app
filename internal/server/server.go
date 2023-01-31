package server

import (
	"context"
	"net/http"
	"tests_app/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(config config.ServerConfig, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		MaxHeaderBytes: config.MaxHeaderBytes,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
