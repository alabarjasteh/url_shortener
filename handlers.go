package main

import (
	"net/http"

	"github.com/alabarjasteh/url-shortener/shortener"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func (c *Controller) CreateShortUrl(w http.ResponseWriter, req *http.Request) {
	paste, err := DecodePostUrl(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := shortener.GenerateShortLink(paste.OriginalLink)

	// write through
	err = c.memcache.Set(shortUrl, paste.OriginalLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.db.Save(shortUrl, paste.OriginalLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := postUrlResponse{ShortLink: shortUrl}
	EncodeResponse(w, response)
	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) HandleShortUrlRedirect(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	shortUrl := params["shortlink"]
	originalUrl, err := c.memcache.Get(shortUrl)
	if err == redis.Nil {
		// does not exist in cache
		// retrive from DB
		originalUrl, err = c.db.Load(shortUrl)

		// write back into cache
		err2 := c.memcache.Set(shortUrl, originalUrl)

		if err2 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, req, originalUrl, http.StatusSeeOther)
	}
}
