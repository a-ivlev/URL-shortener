package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func RedirectPage(w http.ResponseWriter, r *http.Request)  {
	srvHost := os.Getenv("SHORT_SRV_HOST")
	if srvHost == "" {
		srvHost = "localhost"
	}

	//redirectPath := r.URL.Path
	redirectPath := chi.URLParam(r, "short")

	redirect := fmt.Sprintf("%s/%s", srvHost, redirectPath)

	http.Redirect(w, r, redirect, http.StatusFound)
}
