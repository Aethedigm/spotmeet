package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InitRoutes(r chi.Router) {
	r.Route("/test-api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test"))
		})
	})
}
