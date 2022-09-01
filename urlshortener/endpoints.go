package urlshortenersvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type (
	postURLRequest struct {
		OriginalURL string `json:"original_link"`
	}

	postURLResponse struct {
		ShortURL string `json:"short_link"`
		Err      error  `json:"err,omitempty"`
	}

	getURLRequest struct {
		ShortURL string `json:"short_link"`
	}

	getURLResponse struct {
		OriginalURL string
		Err         error `json:"err,omitempty"`
	}
)

// any endpoint that could return an error, implements errorer interface
func (r postURLResponse) error() error { return r.Err }
func (r getURLResponse) error() error  { return r.Err }

// implements redirecter interface
// helps to handle redirecting in trasport layer
func (r getURLResponse) redirect() string { return r.OriginalURL }

// returns RPC-like API around the Shortener Service
func MakePostURLEndpoint(s Shortener, repo Repository, cache Cache) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postURLRequest)
		shortURL, err := s.PostURL(ctx, req.OriginalURL)
		return postURLResponse{ShortURL: shortURL, Err: err}, nil
	}
}

func MakeGetURLEndpoint(s Shortener, repo Repository, cache Cache) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getURLRequest)
		originalURL, err := s.GetURL(ctx, req.ShortURL)
		return getURLResponse{OriginalURL: originalURL, Err: err}, nil
	}
}
