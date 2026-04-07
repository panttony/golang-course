package http

import (
	"net/http"
	"time"

	http_v1 "github.com/panttony/golang-course/api/gen/http/gateway"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, handler http_v1.StrictServerInterface) *Server {
	strictHandler := http_v1.NewStrictHandlerWithOptions(handler, nil, http_v1.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	})

	mux := http.NewServeMux()
	registerSwaggerRoutes(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           http_v1.HandlerFromMux(strictHandler, mux),
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
