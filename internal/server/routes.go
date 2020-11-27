package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"
)

func (s *Server) addRouter() {
	r := chi.NewRouter()
	s.router = r

	r.Use(
		cmw.RequestID,
		cmw.RealIP,
		cmw.Logger,
		cmw.Recoverer,
		cmw.Timeout(5*time.Second),
		cmw.Heartbeat("/ping"),
	)
	if s.useSCS {
		r.Use(s.session.LoadAndSave)
		r.Use(s.tracked.LoadAndSave)
	}

	r.Get("/", Root)
	r.Mount("/saml/", s.saml)
	r.Handle("/secret", s.saml.RequireAccount(http.HandlerFunc(Secret)))
	r.Handle("/profile", s.saml.RequireAccount(http.HandlerFunc(Profile)))
}
