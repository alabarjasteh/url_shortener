package main

import (
	"encoding/json"
	"net/http"

	"github.com/alabarjasteh/url-shortener/shortener"
	"github.com/gorilla/mux"
)

func MakeRoutes(svc shortener.Interface) *mux.Router {
	router := mux.NewRouter()
	router.Methods("POST").Path("/paste").HandlerFunc(svc.PostUrl)
	router.Methods("GET").Path("/{shortlink}").HandlerFunc(svc.RedirectShortUrl)
	return router
}

type postUrlRequest struct {
	OriginalLink string `json:"original_link"`
}

type postUrlResponse struct {
	ShortLink string `json:"short_link"`
}

func DecodePostUrl(r *http.Request) (*postUrlRequest, error) {
	var req postUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
