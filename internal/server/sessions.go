package server

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const sessionCookieName = "sessionid"
const sessionLifetime = 7 * 24 * time.Hour

func (s *Server) addSCS() {
	sm := scs.New()
	sm.Cookie.HttpOnly = true
	sm.Cookie.Name = sessionCookieName
	sm.Cookie.Persist = true
	sm.Cookie.SameSite = http.SameSiteStrictMode
	sm.Cookie.Secure = false // TODO
	sm.IdleTimeout = time.Hour
	sm.Lifetime = sessionLifetime
	s.sm = sm
	s.useSCS = true
}
