package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Shortener, logger log.Logger, repo Repository, cache Cache) http.Handler {
	r := mux.NewRouter()
	r.Methods("POST").Path("/paste").Handler(httptransport.NewServer(
		MakePostURLEndpoint(s, repo, cache),
		decodePostURLRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/{shortlink}").Handler(httptransport.NewServer(
		MakeGetURLEndpoint(s, repo, cache),
		decodeGetURLRequest,
		encodeResponse,
	))
	return r
}

func decodePostURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postURLRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}
func decodeGetURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	shortLink, ok := vars["shortlink"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getURLRequest{ShortURL: shortLink}, nil
}

type (
	errorer interface {
		error() error
	}
	redirecter interface {
		redirect() string
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// business-logic errors
		encodeError(ctx, e.error(), w)
		return nil
	}
	if r, ok := response.(redirecter); ok && r != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Location", r.redirect())
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
