package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alexedwards/scs/v2"
	"github.com/crewjam/saml/samlsp"
	"github.com/go-chi/chi"

	"github.com/sjansen/go-saml-demo/internal/config"
)

var _ samlsp.RequestTracker = &Server{}
var _ samlsp.Session = &Server{}

// Server provides Strongbox's API
type Server struct {
	config     *config.Config
	relaystate *scs.SessionManager
	router     *chi.Mux
	saml       *samlsp.Middleware
	sessions   *scs.SessionManager

	useSCS bool

	done chan struct{}
	wg   sync.WaitGroup
}

// New creates a new Server
func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		config: cfg,
		done:   make(chan struct{}),
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
		cfg, err := config.NewBoltStoreConfig()
		if err != nil {
			return nil, err
		}
		relaystate, sessions, err := s.openBoltStores(cfg)
		if err != nil {
			return nil, err
		}
		s.addSCS(relaystate, sessions)
	case config.DynamoStore:
		cfg, err := config.NewDynamoStoreConfig()
		if err != nil {
			return nil, err
		}
		relaystate, sessions, err := s.openDynamoStores(cfg)
		if err != nil {
			return nil, err
		}
		s.addSCS(relaystate, sessions)
	default:
		return nil, fmt.Errorf("not implemented: %s", cfg.SessionStore)
	}

	s.addRouter()
	return s, nil
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() error {
	fmt.Println("Using session store:", s.config.SessionStore)
	fmt.Println("Listening to", s.config.Root.String())

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.router,
	}
	go func() {
		ch := make(chan os.Signal, 10)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		fmt.Println("Exiting...")
		_ = server.Shutdown(context.Background())
	}()

	err := server.ListenAndServe()
	close(s.done)
	s.wg.Wait()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
