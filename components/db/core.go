package db

import (
	"errors"
	"fmt"

	"github.com/bluest-eel/state/components/config"
)

// db-wide constants
const (
	// Connection types
	BADGER    string = "badger"
	COCKROACH string = "cockroach"
	TESTSTORE string = "teststore"
	UNDEFINED string = "undefined"
)

// DB ...
type DB interface {
	Close() error
	DBType() string
	Get(key string) (*KV, error)
	Set(*KV) error
}

// Open dispatches the creation of a specific DB connection type
func Open(cfg *config.Config) (DB, error) {
	switch cfg.DB.Type {
	case BADGER:
		conn, err := NewBadgerConnector(cfg)
		if err != nil {
			return conn, err
		}
		return conn, nil
	case COCKROACH:
		conn, err := NewCockroachConnector(cfg)
		if err != nil {
			return conn, err
		}
		return conn, nil
	case TESTSTORE:
		return NewTestStoreConnector(cfg)
	default:
		errMsg := fmt.Sprintf("database connection type %s not supported", cfg.DB.Type)
		err := errors.New(errMsg)
		return NewUndefinedConnection(err)
	}
}
