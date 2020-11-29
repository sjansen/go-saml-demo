package server

import (
	"time"

	"github.com/alexedwards/scs/boltstore"
	"github.com/alexedwards/scs/v2"
	"go.etcd.io/bbolt"

	"github.com/sjansen/go-saml-demo/internal/config"
)

func (s *Server) openBoltStores(cfg *config.BoltStoreConfig) (scs.Store, scs.Store, error) {
	db, err := bbolt.Open(cfg.Path+"relaystate.db", 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	relaystate := boltstore.NewWithCleanupInterval(db, time.Minute)
	go func(db *bbolt.DB) {
		<-s.done
		db.Close()
		s.wg.Done()
	}(db)
	s.wg.Add(1)

	db, err = bbolt.Open(cfg.Path+"sessions.db", 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	sessions := boltstore.NewWithCleanupInterval(db, time.Minute)
	go func(db *bbolt.DB) {
		<-s.done
		db.Close()
		s.wg.Done()
	}(db)
	s.wg.Add(1)

	return relaystate, sessions, nil
}
