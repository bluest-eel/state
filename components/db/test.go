package db

import (
	"time"

	"github.com/bluest-eel/state/components/config"
)

// Constants for testing
const (
	TestCreatedTime = "2006-01-02T15:04:05Z"
	TestUpdatedTime = "2006-01-08T10:00:08Z"
	TestKey         = "a key"
	TestVal         = "a value"
)

// TestStoreConnector is for test dbs that will be used in unit tests
type TestStoreConnector struct {
	ConnType string
}

// NewTestStoreConnector is for use with unit tests
func NewTestStoreConnector(_ *config.Config) (*TestStoreConnector, error) {
	connector := &TestStoreConnector{ConnType: TESTSTORE}
	return connector, nil
}

// Close closes the connection to the database
func (c *TestStoreConnector) Close() error {
	return nil
}

// DBType returns the database type
func (c *TestStoreConnector) DBType() string {
	return c.ConnType
}

// Get ...
func (c *TestStoreConnector) Get(key string) (*KV, error) {
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
func (c *TestStoreConnector) Set(inpuc *KV) error {
	return nil
}
