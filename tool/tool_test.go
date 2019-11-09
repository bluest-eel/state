package tool_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/bluest-eel/state/components"
	"github.com/bluest-eel/state/tool"
	"github.com/stretchr/testify/suite"
)

type toolTestSuite struct {
	components.TestBase
}

func TestToolTestSuite(t *testing.T) {
	suite.Run(t, &toolTestSuite{})
}

func (suite *toolTestSuite) TestParsePoints() {
	result := tool.ParsePoints(tool.OxfordPubs)
	names := make([]string, len(result))
	for i, p := range result {
		names[i] = p.Name
	}
	sort.Strings(names)
	suite.Equal(
		"Morse Bar,The Eagle and Child,The King's Arms,The Turf Tavern,The White Horse",
		strings.Join(names, ","))
}

func (suite *toolTestSuite) TestSplitter1Arg() {
	var latArg, lonArg, name string
	err := tool.Splitter("51.770903", tool.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("", lonArg)
	suite.Equal("", name)
	err = tool.Splitter("51.770903::", tool.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("", lonArg)
	suite.Equal("", name)
}

func (suite *toolTestSuite) TestSplitter2Args() {
	var latArg, lonArg, name string
	err := tool.Splitter("51.770903::-1.2626219", tool.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("-1.2626219", lonArg)
	suite.Equal("", name)
}

func (suite *toolTestSuite) TestSplitter3Args() {
	var latArg, lonArg, name string
	err := tool.Splitter(tool.TolkiensHouse, tool.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("-1.2626219", lonArg)
	suite.Equal("Tolkien's House", name)
}

func (suite *toolTestSuite) TestSplitterMoreThan3Args() {
	var latArg, lonArg, name string
	err := tool.Splitter(tool.TolkiensHouse+"::extraneous", tool.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.Error(err)
	suite.Equal("Too many delimited items", err.Error())
	err = tool.Splitter(tool.TolkiensHouse+"::extra::neo::us", tool.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.Error(err)
	suite.Equal("Too many delimited items", err.Error())
}
