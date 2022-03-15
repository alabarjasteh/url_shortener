package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Controller) MakeRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Methods("POST").Path("/paste").HandlerFunc(c.CreateShortUrl)
	router.Methods("GET").Path("/{shortlink}").HandlerFunc(c.HandleShortUrlRedirect)
	return router
}

type postUrlRequest struct {
	OriginalLink string `json:"originalLink"`
}

type postUrlResponse struct {
	ShortLink string `json:"shortLink"`
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
