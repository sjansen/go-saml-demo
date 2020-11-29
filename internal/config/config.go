package config

import (
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

// BoltStoreConfig contains settings required for BoltStore
type BoltStoreConfig struct {
	Path string `envconfig:"GSD_SESSION_DIR"`
}

// DynamoStoreConfig contains settings required for DynamoStore
type DynamoStoreConfig struct {
	TableName string `envconfig:"GSD_SESSION_TABLE"`
}

// New loads application settings from the environment.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewBoltStoreConfig loads BoltStore settings from the environment.
func NewBoltStoreConfig() (*BoltStoreConfig, error) {
	cfg := &BoltStoreConfig{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewDynamoStoreConfig loads DynamoStore settings from the environment.
func NewDynamoStoreConfig() (*DynamoStoreConfig, error) {
	cfg := &DynamoStoreConfig{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
