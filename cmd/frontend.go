package main

import (
	"github.com/a-ivlev/URL-shortener/internal/handlers"
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

	cliPort := os.Getenv("PORT")
	if cliPort == "" {
		log.Fatal("unknown PORT = ", cliPort)
	}

	err := http.ListenAndServe(":"+cliPort, r)
	if err != nil {
		log.Println("Client shortener stopped...")
	}
}
