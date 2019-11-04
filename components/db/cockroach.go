package db

import (
	"time"

	"github.com/bluest-eel/state/components/config"
)

// CockroachConnector ...
type CockroachConnector struct {
	ConnType string
	// Conn
}

// NewCockroachConnector ...
func NewCockroachConnector(_ *config.Config) (*CockroachConnector, error) {
	connector := &CockroachConnector{ConnType: TESTSTORE}
	return connector, nil
}

// Close closes the connection to the database
func (c *CockroachConnector) Close() error {
	return nil
}

// DBType returns the database type
func (c *CockroachConnector) DBType() string {
	return c.ConnType
}

// Get ...
func (c *CockroachConnector) Get(key string) (*KV, error) {
	createdOn, _ := time.Parse(time.RFC3339, TestCreatedTime)
	updatedOn, _ := time.Parse(time.RFC3339, TestUpdatedTime)
	return &KV{
		Key:     key,
		Value:   []byte(TestVal),
		Created: &createdOn,
		Updated: &updatedOn,
	}, nil
}

// Set ...
func (c *CockroachConnector) Set(inpuc *KV) error {
	return nil
}
