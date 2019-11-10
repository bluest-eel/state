package nearest_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/bluest-eel/state/components"
	"github.com/bluest-eel/state/examples/nearest"
	"github.com/golang/geo/s1"
	"github.com/stretchr/testify/suite"
)

type nearestTestSuite struct {
	components.TestBase
}

func TestNearestTestSuite(t *testing.T) {
	suite.Run(t, &nearestTestSuite{})
}

func (suite *nearestTestSuite) TestParsePoint() {
	result := nearest.ParsePoint("51.7546135::-1.2577909::The White Horse")
	suite.Equal("4876c6af5db009dd", result.CellID.ToToken())
	suite.Equal(s1.Angle(0.9032884086721062), result.LatLon.Lat)
	suite.Equal(s1.Angle(-0.021952592506622747), result.LatLon.Lng)
	suite.Equal(51.7546135, result.LatLon.Lat.Degrees())
	suite.Equal(-1.2577909, result.LatLon.Lng.Degrees())
}

func (suite *nearestTestSuite) TestParsePoints() {
	result := nearest.ParsePoints(nearest.OxfordPubExamplePoints)
	names := make([]string, len(result))
	for i, p := range result {
		names[i] = p.Name
	}
	sort.Strings(names)
	suite.Equal(
		"Morse Bar,The Eagle and Child,The King's Arms,The Turf Tavern,The White Horse",
		strings.Join(names, ","))
}

func (suite *nearestTestSuite) TestSplitter1Arg() {
	var latArg, lonArg, name string
	err := nearest.Splitter("51.770903", nearest.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("", lonArg)
	suite.Equal("", name)
	err = nearest.Splitter("51.770903::", nearest.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("", lonArg)
	suite.Equal("", name)
}

func (suite *nearestTestSuite) TestSplitter2Args() {
	var latArg, lonArg, name string
	err := nearest.Splitter("51.770903::-1.2626219", nearest.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("-1.2626219", lonArg)
	suite.Equal("", name)
}

func (suite *nearestTestSuite) TestSplitter3Args() {
	var latArg, lonArg, name string
	err := nearest.Splitter(nearest.TolkiensHouse, nearest.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.NoError(err)
	suite.Equal("51.770903", latArg)
	suite.Equal("-1.2626219", lonArg)
	suite.Equal("Tolkien's House", name)
}

func (suite *nearestTestSuite) TestSplitterMoreThan3Args() {
	var latArg, lonArg, name string
	err := nearest.Splitter(nearest.TolkiensHouse+"::extraneous", nearest.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.Error(err)
	suite.Equal("Too many delimited items", err.Error())
	err = nearest.Splitter(nearest.TolkiensHouse+"::extra::neo::us", nearest.PointSplitterOpts,
		&latArg, &lonArg, &name)
	suite.Error(err)
	suite.Equal("Too many delimited items", err.Error())
}
