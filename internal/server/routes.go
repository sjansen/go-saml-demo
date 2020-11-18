package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/go-saml-demo/internal/config"
)

func (s *Server) addRouter(cfg *config.Config) error {
	r := chi.NewRouter()
	s.router = r

	samlSP, err := newSAMLMiddleware(cfg)
	if err != nil {
		return err
	}

	r.Use(
		cmw.RequestID,
		cmw.RealIP,
		cmw.Logger,
		cmw.Recoverer,
		cmw.Timeout(5*time.Second),
		cmw.Heartbeat("/ping"),
		s.sm.LoadAndSave,
	)

	r.Get("/", Root)
	r.Mount("/saml/", samlSP)
	r.Handle("/secret", samlSP.RequireAccount(http.HandlerFunc(Secret)))

	return nil
}
