package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/server/http/handlers"
)

type Server struct {
	logger repository.Logger
	app    repository.Application
	server *http.Server
}

func NewServer(host, port string, logger repository.Logger, app repository.Application) *Server {
	mux := http.NewServeMux()
	httpHandlers := handlers.NewHTTPHandler(logger)
	mux.HandleFunc("/", httpHandlers.HomeHandler)
	middleware := NewMiddleware(logger)
	handler := middleware.loggingMiddleware(mux)
	server := http.Server{
		Addr:         net.JoinHostPort(host, port),
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	return &Server{
		logger: logger,
		app:    app,
		server: &server,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting http server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping http server...")
	return s.server.Shutdown(ctx)
}
