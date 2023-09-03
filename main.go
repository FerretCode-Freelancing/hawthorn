package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

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

	err := os.MkdirAll("/tmp/hawthorn", os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("/tmp/hawthorn/cache.json"); os.IsNotExist(err) == true {
		err = os.WriteFile("/tmp/hawthorn/cache.json", []byte("{}"), fs.ModePerm)

		if err != nil {
			log.Fatal(err)
		}
	}


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

			http.Error(w, "There was an error logging you in.", http.StatusInternalServerError)
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

			http.Error(w, "There was an error listing all containers.", http.StatusInternalServerError)
		}
	})

	r.Post("/containers/new", func(w http.ResponseWriter, r *http.Request) {
		err := containers.New(w, r, o)

		if err != nil {
			fmt.Println(err)

			http.Error(w, "There was an error creating a new container.", http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":3006", r)
}
