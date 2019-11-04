package db

import (
	"github.com/bluest-eel/common/util"
	"github.com/bluest-eel/state/components/config"
	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

// BadgerConnector ...
type BadgerConnector struct {
	ConnType string
	Conn     *badger.DB
}

// NewBadgerConnector ...
func NewBadgerConnector(cfg *config.Config) (*BadgerConnector, error) {
	conn, err := badger.Open(badger.DefaultOptions(cfg.DB.Badger.Directory))
	if err != nil {
		log.Fatal(err)
	}
	connector := &BadgerConnector{
		ConnType: BADGER,
		Conn:     conn,
	}

	return connector, err
}

// Close closes the connection to the database
func (c *BadgerConnector) Close() error {
	return c.Conn.Close()
}

// DBType returns the database type
func (c *BadgerConnector) DBType() string {
	return c.ConnType
}

// Get ...
func (c *BadgerConnector) Get(key string) (*KV, error) {
	decoded := &KV{}
	err := c.Conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = util.GOBDecode(valCopy, decoded)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return decoded, err
}

// Set ...
func (c *BadgerConnector) Set(input *KV) error {
	err := c.Conn.Update(func(txn *badger.Txn) error {
		encoded, err := util.GOBEncode(input)
		if err != nil {
			log.Error(err)
		}
		err = txn.Set([]byte(input.Key), encoded)
		if err != nil {
			log.Error(err)
		}
		return err
	})
	return err
}
