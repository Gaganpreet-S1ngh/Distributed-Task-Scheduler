package server

import (
	"net/http"
	"time"
)

func NewHTTPServer(port string) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
