package main

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func (s *ShortenerService) PostUrl(w http.ResponseWriter, req *http.Request) {
	paste, err := DecodePostUrl(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := GenerateShortLink(paste.OriginalLink)

	// write through
	err = s.memcache.Set(shortUrl, paste.OriginalLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.Save(shortUrl, paste.OriginalLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := postUrlResponse{ShortLink: shortUrl}
	EncodeResponse(w, response)
	w.WriteHeader(http.StatusCreated)
}

func (s *ShortenerService) RedirectShortUrl(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	shortUrl := params["shortlink"]
	originalUrl, err := s.memcache.Get(shortUrl)
	if err == redis.Nil {
		// does not exist in cache
		// retrive from DB
		originalUrl, err = s.db.Load(shortUrl)

		// write back into cache
		err2 := s.memcache.Set(shortUrl, originalUrl)

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
