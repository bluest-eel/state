package tool_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/bluest-eel/state/components"
	"github.com/bluest-eel/state/tool"
	"github.com/golang/geo/s1"
	"github.com/stretchr/testify/suite"
)

type toolTestSuite struct {
	components.TestBase
}

func TestToolTestSuite(t *testing.T) {
	suite.Run(t, &toolTestSuite{})
}

func (suite *toolTestSuite) TestParsePoint() {
	result := tool.ParsePoint("51.7546135::-1.2577909::The White Horse")
	suite.Equal("4876c6af5db009dd", result.CellID.ToToken())
	suite.Equal(s1.Angle(0.9032884086721062), result.LatLon.Lat)
	suite.Equal(s1.Angle(-0.021952592506622747), result.LatLon.Lng)
	suite.Equal(51.7546135, result.LatLon.Lat.Degrees())
	suite.Equal(-1.2577909, result.LatLon.Lng.Degrees())
}

func (suite *toolTestSuite) TestParsePoints() {
	result := tool.ParsePoints(tool.OxfordPubExamplePoints)
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
