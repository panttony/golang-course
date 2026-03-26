package http

import (
	"net/http"
	"time"

	gen "github.com/pantonny/golang-course/internal/gateway/transport/http/gen"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, handler gen.StrictServerInterface) *Server {
	strictHandler := gen.NewStrictHandler(handler, nil)
	mux := http.NewServeMux()
	gen.HandlerFromMux(strictHandler, mux)

	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
