package marshaller_test

import (
	"testing"

	"github.com/bluest-eel/state/components"
	"github.com/bluest-eel/state/components/config"
	"github.com/bluest-eel/state/components/marshaller"
	"github.com/stretchr/testify/suite"
)

type gobTestSuite struct {
	components.TestBase
	cfg        *config.Config
	marshaller marshaller.Marsh
}

func (suite *gobTestSuite) SetupSuite() {
	suite.cfg = config.NewConfig()
	suite.cfg.Marshaller.Type = marshaller.GOB
	suite.marshaller, _ = marshaller.New(suite.cfg)
}

func TestGobTestSuite(t *testing.T) {
	suite.Run(t, &gobTestSuite{})
}

func (suite *gobTestSuite) TestMarshalUnmarhsal() {
	testData := &marshaller.StateMetadata{Name: "a"}
	marshalled, err := suite.marshaller.Marshal(testData)
	suite.NoError(err)
	unmarshalled, err := suite.marshaller.Unmarshal(marshalled)
	suite.NoError(err)
	suite.Equal("XXX", unmarshalled)
	// suite.Equal("a", unmarshalled.Thing1)
	// suite.Equal("a", unmarshalled.Thing2)
}
