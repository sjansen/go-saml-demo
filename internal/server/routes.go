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
		s.sm.LoadAndSave,
	)

	r.Get("/", Root)
	r.Mount("/saml/", s.sp)
	r.Handle("/secret", s.sp.RequireAccount(http.HandlerFunc(Secret)))
}
