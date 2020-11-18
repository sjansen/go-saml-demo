package config

import "github.com/vrischmann/envconfig"

// Config contains application settings.
type Config struct {
	Debug bool `envconfig:"DEBUG,optional"`
	SAML  struct {
		EntityID    string `envconfig:"SAML_ENTITY_ID,default=go-saml-demo"`
		RootURL     string `envconfig:"SAML_BASE_URL,default=http://localhost:8080/"`
		MetadataURL string `envconfig:"SAML_METADATA_URL"`
		Certificate string `envconfig:"SAML_CERTIFICATE"`
		PrivateKey  string `envconfig:"SAML_PRIVATE_KEY"`
	}
}

// New loads application settings from the environment.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
