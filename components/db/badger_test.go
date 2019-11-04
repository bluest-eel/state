package db_test

import (
	"testing"

	"github.com/bluest-eel/state/components"
	"github.com/bluest-eel/state/components/config"
	"github.com/bluest-eel/state/components/db"
	"github.com/stretchr/testify/suite"
)

type badgerTestSuite struct {
	components.TestBase
	cfg *config.Config
	db  db.DB
}

func (suite *badgerTestSuite) SetupSuite() {
	suite.cfg = config.NewConfig()
	suite.cfg.DB.Type = db.BADGER
	suite.cfg.DB.Badger.Directory = "/tmp/testdif"
	suite.db, _ = db.NewBadgerConnector(suite.cfg)
}

func TestBadgerTestSuite(t *testing.T) {
	suite.Run(t, &badgerTestSuite{})
}

func (suite *badgerTestSuite) TestSetGet() {
	testData := &db.KV{Key: "stuff", Value: []byte("things")}
	_ = suite.db.Set(testData)
	result, _ := suite.db.Get("stuff")
	suite.Equal("things", string(result.Value))
}
