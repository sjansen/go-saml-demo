package server

import (
	"time"

	"github.com/alexedwards/scs/boltstore"
	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.etcd.io/bbolt"

	"github.com/sjansen/dynamostore"
	"github.com/sjansen/go-saml-demo/internal/config"
)

func (s *Server) openBoltStores(cfg *config.BoltStoreConfig) (scs.Store, scs.Store, error) {
	db, err := bbolt.Open(cfg.Path+"sessions.db", 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	store := boltstore.NewWithCleanupInterval(db, time.Minute)
	go func(db *bbolt.DB) {
		<-s.done
		db.Close()
		s.wg.Done()
	}(db)
	s.wg.Add(1)

	relaystate := NewPrefixStore("r:", store)
	sessions := NewPrefixStore("s:", store)
	return relaystate, sessions, nil
}

func (s *Server) openDynamoStores(cfg *config.DynamoStoreConfig) (scs.Store, scs.Store, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, nil, err
	}

	awscfg := aws.NewConfig().
		WithCredentials(
			credentials.NewStaticCredentials("id", "secret", "token"),
		)
	if cfg.Endpoint.Host != "" {
		awscfg = awscfg.WithEndpoint(cfg.Endpoint.String())
	}
	svc := dynamodb.New(sess, awscfg)

	store := dynamostore.NewWithTableName(svc, cfg.Table)
	if cfg.Create {
		err := store.CreateTable()
		if err != nil {
			return nil, nil, err
		}
	}

	relaystate := NewPrefixStore("r:", store)
	sessions := NewPrefixStore("s:", store)
	return relaystate, sessions, nil
}

// PrefixStore enables multiple sessions to be stored in a single
// session store by automatically pre-pending a prefix to tokens.
type PrefixStore struct {
	prefix string
	store  scs.Store
}

// NewPrefixStore wraps a session store so it can be shared.
func NewPrefixStore(prefix string, store scs.Store) *PrefixStore {
	return &PrefixStore{
		prefix: prefix,
		store:  store,
	}
}

// Delete removes the session token and data from the store.
func (s *PrefixStore) Delete(token string) (err error) {
	return s.store.Delete(s.prefix + token)
}

// Find returns the data for a session token from the store.
func (s *PrefixStore) Find(token string) (b []byte, found bool, err error) {
	return s.store.Find(s.prefix + token)
}

// Commit add the session token and data to the store.
func (s *PrefixStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	return s.store.Commit(s.prefix+token, b, expiry)
}
