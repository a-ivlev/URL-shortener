package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func RedirectPage(w http.ResponseWriter, r *http.Request)  {
	srvHost := os.Getenv("SHORT_SRV_HOST")
	if srvHost == "" {
		srvHost = "localhost"
	}

	srvPort := os.Getenv("SHORT_SRV_PORT")
	if srvPort == "" {
		srvPort = "8035"
	}
	redirectPath := r.URL.Path
	//redirectPath := chi.URLParam(r, "short")

	redirect := fmt.Sprintf("http://%s:%s%s", srvHost, srvPort, redirectPath)

	http.Redirect(w, r, redirect, http.StatusFound)
}
