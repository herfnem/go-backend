package routes

import (
	"net/http"

	"github.com/swaggo/http-swagger"
)

func RegisterSwaggerRoutes(mux *http.ServeMux) {
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
}
