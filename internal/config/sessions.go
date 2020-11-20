package config

import "fmt"

// BoltStoreConfig contains settings required for BoltStore
type BoltStoreConfig struct {
	Path string `envconfig:"GSD_SESSION_FILE"`
}

// DynamoStoreConfig contains settings required for DynamoStore
type DynamoStoreConfig struct {
	TableName string `envconfig:"GSD_SESSION_TABLE"`
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
	case "boltdb":
		*store = BoltStore
	case "dynamodb":
		*store = DynamoStore
	default:
		return fmt.Errorf("invalid session store: %q", s)
	}
	return nil
}

func (store SessionStore) String() string {
	switch store {
	case DefaultStore:
		return "Default"
	case BoltStore:
		return "BoltStore"
	case DynamoStore:
		return "DynamoStore"
	default:
		return "Invalid"
	}
}
