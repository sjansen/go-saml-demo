package server

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/crewjam/saml/samlsp"
	"github.com/go-chi/chi"

	"github.com/sjansen/go-saml-demo/internal/config"
)

var _ samlsp.RequestTracker = &Server{}
var _ samlsp.Session = &Server{}

// Server provides Strongbox's API
type Server struct {
	useSCS bool

	config  *config.Config
	router  *chi.Mux
	saml    *samlsp.Middleware
	session *scs.SessionManager
	tracked *scs.SessionManager
}

// New creates a new Server
func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		config: cfg,
	}

	sp, err := newSAMLMiddleware(cfg)
	if err != nil {
		return nil, err
	}
	s.saml = sp

	switch cfg.SessionStore {
	case config.DefaultStore:
		// noop
	case config.BoltStore:
		if _, err := config.NewBoltStoreConfig(); err != nil {
			return nil, err
		}
		s.addSCS()
	case config.DynamoStore:
		if _, err := config.NewDynamoStoreConfig(); err != nil {
			return nil, err
		}
		s.addSCS()
	default:
		return nil, fmt.Errorf("not implemented: %s", cfg.SessionStore)
	}

	s.addRouter()
	return s, nil
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe(addr string) error {
	fmt.Println("Using session store:", s.config.SessionStore)
	fmt.Println("Listening to", addr)
	return http.ListenAndServe(addr, s.router)
}
