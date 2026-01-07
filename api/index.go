package handler

import (
	"net/http"
	"github.com/HCMUT-LostFound/backend"
)

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	backend.Handler(w, r)
}

