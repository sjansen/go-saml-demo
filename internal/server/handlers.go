package server

import (
	"html/template"
	"net/http"

	"github.com/crewjam/saml/samlsp"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").ParseGlob("templates/*.html"))
}

// Profile shows user attributes when SCS isn't enabled.
func Profile(w http.ResponseWriter, r *http.Request) {
	var attrs samlsp.Attributes
	s := samlsp.SessionFromContext(r.Context())
	if sa, ok := s.(samlsp.SessionWithAttributes); ok {
		attrs = sa.GetAttributes()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "profile.html", attrs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Root is the root app page.
func Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "welcome.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Secret is a protected page.
func Secret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "secret.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
