package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ferretcode-freelancing/hawthorn/orchestrator"
	"github.com/ferretcode-freelancing/hawthorn/routes/auth"
	"github.com/ferretcode-freelancing/hawthorn/routes/containers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	ctx := context.Background()

	o, err := orchestrator.NewOrchestrator(orchestrator.Orchestrator{
		Context: ctx,
	})

	if err != nil {
		log.Fatal(err)
	}

	// auth
	r.Get("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		err := auth.Login(w, r)

		if err != nil {
			fmt.Println(err)
		}
	})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		err := auth.Callback(w, r)

		if err != nil {
			fmt.Println(err)
		}
	})

	// containers
	r.Get("/containers/list", func(w http.ResponseWriter, r *http.Request) {
		err := containers.List(w, r, o)

		if err != nil {
			fmt.Println(err)
		}
	})

	r.Post("/containers/new", func(w http.ResponseWriter, r *http.Request) {
		err := containers.New(w, r, o)

		if err != nil {
			fmt.Println(err)
		}
	})

	http.ListenAndServe(":3006", r)
}
