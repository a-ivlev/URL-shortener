package main

import (
	"CourseProjectBackendDevGoLevel-1/client/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", handlers.HomePage)
	r.Post("/", handlers.HomePage)
	r.Get("/{short}", handlers.RedirectPage)
	r.Get("/stat/{stat}", handlers.StatPage)

	cliPort := os.Getenv("SHORT_CLI_PORT")
	if cliPort == "" {
		cliPort = "8080"
	}

	err := http.ListenAndServe(":"+cliPort, r)
	if err != nil {
		log.Println("Client shortener stopped...")
	}
}
