package db

import "time"

// UndefinedConnector is for situations where an implementation needs to be
// returned, even when there is none.
type UndefinedConnector struct {
	ConnType string
}

// NewUndefinedConnection is for use with unit tests
func NewUndefinedConnection(err error) (*UndefinedConnector, error) {
	connector := &UndefinedConnector{ConnType: UNDEFINED}
	return connector, err
}

// Close closes the connection to the database
func (c *UndefinedConnector) Close() error {
	return nil
}

// DBType returns the database type
func (c *UndefinedConnector) DBType() string {
	return c.ConnType
}

// Get ...
func (c *UndefinedConnector) Get(key string) (*KV, error) {
	createdOn, _ := time.Parse(time.RFC3339, TestCreatedTime)
	updatedOn, _ := time.Parse(time.RFC3339, TestUpdatedTime)
	return &KV{
		Key:     "",
		Value:   []byte(""),
		Created: &createdOn,
		Updated: &updatedOn,
	}, nil
}

// Set ...
func (c *UndefinedConnector) Set(inpuc *KV) error {
	return nil
}
