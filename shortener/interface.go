package shortener

import "net/http"

type Interface interface {
	PostUrl(w http.ResponseWriter, req *http.Request)
	RedirectShortUrl(w http.ResponseWriter, req *http.Request)
}
