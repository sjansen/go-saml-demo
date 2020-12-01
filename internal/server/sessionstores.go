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

	relaystate := dynamostore.NewWithTableName(svc, cfg.Table)
	sessions := dynamostore.NewWithTableName(svc, cfg.Table)
	if cfg.Create {
		err := sessions.CreateTable()
		if err != nil {
			return nil, nil, err
		}
	}

	return relaystate, sessions, nil
}
