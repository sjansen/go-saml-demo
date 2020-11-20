package config

import (
	"fmt"
	"net/url"

	"github.com/vrischmann/envconfig"
)

// Config contains application settings.
type Config struct {
	Debug   bool `envconfig:"DEBUG,optional"`
	RootURL URL  `envconfig:"GSD_URL,default=http://localhost:8080/"`
	SAML    struct {
		EntityID    string `envconfig:"GSD_SAML_ENTITY_ID,default=go-saml-demo"`
		MetadataURL string `envconfig:"GSD_SAML_METADATA_URL"`
		Certificate string `envconfig:"GSD_SAML_CERTIFICATE"`
		PrivateKey  string `envconfig:"GSD_SAML_PRIVATE_KEY"`
	}
	SessionStore SessionStore `envconfig:"GSD_SESSION_STORE,default=default"`
}

// New loads application settings from the environment.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// SessionStore is an enum of possible session stores
type SessionStore int

const (
	// DefaultStore is cookie-based with only SAML attributes
	DefaultStore SessionStore = iota
	// BoltStore is backed by a local file
	BoltStore
	// DynamoStore is backed by DynamoDB
	DynamoStore
)

// Unmarshal converts an environment variable string to a URL
func (store *SessionStore) Unmarshal(s string) error {
	switch s {
	case "default":
		*store = DefaultStore
	case "bolt":
		*store = BoltStore
	case "dynamodb":
		*store = DynamoStore
	default:
		return fmt.Errorf("invalid session store: %q", s)
	}
	return nil
}

// URL represents a parsed URL
type URL struct {
	url.URL
}

// Unmarshal converts an environment variable string to a URL
func (u *URL) Unmarshal(s string) error {
	return u.URL.UnmarshalBinary([]byte(s))
}
