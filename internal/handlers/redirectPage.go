package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func RedirectPage(w http.ResponseWriter, r *http.Request)  {
	srvHost := os.Getenv("SRV_HOST")
	if srvHost == "" {
		log.Fatal("unknown SRV_HOST = ", srvHost)
	}
	//srvPort := os.Getenv("SRV_PORT")
	//if srvPort == "" {
	//	log.Fatal("unknown SRV_PORT = ", srvPort)
	//}

	redirectPath := chi.URLParam(r, "short")

	//redirect := fmt.Sprintf("http://%s:%s/%s", srvHost, srvPort, redirectPath)
	redirect := fmt.Sprintf("http://%s/%s", srvHost, redirectPath)

	http.Redirect(w, r, redirect, http.StatusFound)
}
