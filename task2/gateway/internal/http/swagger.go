package http

import (
	_ "embed"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const swaggerDocPath = "/swagger/doc.yaml"

//go:embed docs/gateway.yaml
var swaggerSpec []byte

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("GET "+swaggerDocPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(swaggerSpec)
	})

	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL(swaggerDocPath),
	))
}
