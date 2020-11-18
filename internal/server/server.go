package server

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/sjansen/go-saml-demo/internal/config"
)

// Server provides Strongbox's API
type Server struct {
	router *chi.Mux
	sm     *scs.SessionManager
}

// New creates a new Server
func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		sm: newSessionManager(),
	}

	err := s.addRouter(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe(addr string) error {
	fmt.Println("Listening to", addr)
	return http.ListenAndServe(addr, s.router)
}
