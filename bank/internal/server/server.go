package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	port   string
	srv    *http.Server
}

func NewServer(router *mux.Router, port string) *Server {
	return &Server{
		router: router,
		port:   port,
	}
}

func (s *Server) Start() error {
	s.srv = &http.Server{
		Handler:      s.router,
		Addr:         s.port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.srv.ListenAndServe()
}
