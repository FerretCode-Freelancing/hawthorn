package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	// auth
	r.Get("/auth/login", func(w http.ResponseWriter, r *http.Request) {})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {})

	// containers
	r.Get("/containers/list", func(w http.ResponseWriter, r *http.Request) {})	

	r.Post("/containers/new", func(w http.ResponseWriter, r *http.Request) {})
}