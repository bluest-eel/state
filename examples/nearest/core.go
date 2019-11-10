// Package tool copied from https://gist.github.com/antoniomo/3371e44cbe2f0cc75a525aac0d188cfb
package nearest

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
)

// https://blog.nobugware.com/post/2016/geo_db_s2_geohash_database/
// http://s2geometry.io/devguide/cpp/quickstart

// Tool constants
const (
	EarthRadiusKm   = 6371.01 // Earth's mean radius, km (according to NASA).
	Helsinki        = "helsinki"
	Oxford          = "oxford"
	OxfordPubs      = "oxford-pubs"
	PointDelimiter  = "::"
	PointsDelimiter = "|"
	Tolkien         = "tolkien"
)

// Example Points
var (
	HelsinkiCenter        = "60.1699::24.9384::Helsinki Center"
	HelsinkiExamplePoints = "60.2934::25.0378::Vantaa Center|" +
		"60.2055::24.6559::Espoo Center|" +
		"60.1699::24.9380::Person in Helsinki|" +
		// Far points
		"50.0::150.0::Seah of Okhotsk|" +
		"50.0::-150.0::North Pacific Ocean|" +
		"50.0::50.0::Western Kazakhstan|" +
		"50.0::-50.0::Labrador Sea"
	OxfordExamplePoints = "51.751944::-1.257778::Oxford|" +
		"51.7572::-1.2603::The Eagle and Child|" +
		"51.507222::-0.1275::London|" +
		"51.48::0::Greenwich|" +
		"52.205278::0.119167::Cambridge|" +
		// Far points
		"42.116667::-71.865278::Oxford, Massachusetts|" +
		"-43.312778::172.190556::Oxford, New Zealand"
	// Pints points
	OxfordPubExamplePoints = "51.7572::-1.2603::The Eagle and Child|" +
		"51.7550609::-1.2617064::Morse Bar|" +
		"51.755::-1.2544::The King's Arms|" +
		"51.7546135::-1.2577909::The White Horse|" +
		"51.7547::-1.253::The Turf Tavern"
	TolkiensHouse = "51.770903::-1.2626219::Tolkien's House"
)

/////////////////////////////////////////////////////////////////////////////
///   Tool Entrypouint   /////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// RunOptions ...
type RunOptions struct {
	CapRadius   float64
	CenterPoint string
	MaxCells    int
	MaxLevel    int
	Points      string
}

// Run https://godoc.org/github.com/golang/geo/s2#Cap
func Run(opts *RunOptions) {
	center := ParsePoint(opts.CenterPoint)
	log.Infof("Center cell id: %#v", center)
	points := ParsePoints(opts.Points)
	points = append(points, center)

	s2cap := s2.CapFromCenterAngle(center.CellID.Point(), KmToAngle(opts.CapRadius))
	// http://s2geometry.io/resources/s2cell_statistics.html
	// Level 12 are about 3 to 6.4km^2 cells
	// Level 20 are about 46.41 to 97.3 meter cells
	// Since we put a MaxCells of 8, it won't go to the max level if it
	// can't approximate the region better anyway.
	rc := &s2.RegionCoverer{MaxLevel: opts.MaxLevel, MaxCells: opts.MaxCells}
	covering := rc.Covering(s2.Region(s2cap))

	for i, cov := range covering {
		log.Infof("Covering Cell %d ID: %d Level: %d", i, uint64(cov),
			cov.Level())
		PointsInCellID(s2cap, cov, center, points)
	}
}

/////////////////////////////////////////////////////////////////////////////
///   Point Implementation   ////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// Point ...
type Point struct {
	CellID s2.CellID
	LatLon s2.LatLng
	Name   string
}

// NewPoint ...
func NewPoint(lat, lon float64, name string) Point {
	latLon := s2.LatLngFromDegrees(lat, lon)
	return Point{
		CellID: s2.CellIDFromLatLng(latLon),
		LatLon: latLon,
		Name:   name,
	}
}

/////////////////////////////////////////////////////////////////////////////
///   Points Functions   ////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// PointsInCellID ...
func PointsInCellID(s2cap s2.Cap, cov s2.CellID, center Point, points []Point) {
	bmin := uint64(cov.RangeMin())
	bmax := uint64(cov.RangeMax())

	for _, p := range points {
		// This simulates an indexed range query on the DB
		if uint64(p.CellID) < bmin || uint64(p.CellID) > bmax {
			continue
		}
		// Only those in range
		lat := p.LatLon.Lat.Degrees()
		lon := p.LatLon.Lng.Degrees()
		log.Infof("Nearby Candidate: %f (lat) %f (lon) %s", lat, lon, p.Name)
		log.Infof("Calculated distance to %s: %f (km)", center.Name,
			AngleToKm(p.LatLon.Distance(center.LatLon)))
		log.Info("False positive? ", !s2cap.ContainsPoint(p.CellID.Point()))
	}
}

/////////////////////////////////////////////////////////////////////////////
///   CLI Helper Functions   ////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// ParsePoint ...
func ParsePoint(delimited string) Point {
	log.Debugf("Pasring point '%s' ...", delimited)
	var latArg, lonArg, name string
	err := Splitter(delimited, PointSplitterOpts, &latArg, &lonArg, &name)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case latArg == Helsinki:
		return ParsePoint(HelsinkiCenter)
	case latArg == Tolkien:
		return ParsePoint(TolkiensHouse)
	}
	lat, lon := ConvertLatLon(latArg, lonArg)
	return NewPoint(lat, lon, name)
}

// ParsePoints ...
func ParsePoints(delimited string) []Point {
	log.Debugf("Pasring points '%s' ...", delimited)
	args := strings.Split(delimited, PointsDelimiter)
	alias := args[0]
	switch {
	case alias == Helsinki:
		return ParsePoints(HelsinkiExamplePoints)
	case alias == Oxford:
		return ParsePoints(OxfordExamplePoints)
	case alias == OxfordPubs:
		return ParsePoints(OxfordPubExamplePoints)
	default:
		points := make([]Point, len(args))
		for i, p := range args {
			points[i] = ParsePoint(p)
		}
		return points
	}
}

/////////////////////////////////////////////////////////////////////////////
///   Utility Functions   ///////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

// KmToAngle converts a distance on the Earth's surface to an angle.
// https://github.com/golang/geo/blob/23949e136d58aeb8aa39844a312b68d90c4eb8aa/s2/s2_test.go#L38-L43
func KmToAngle(km float64) s1.Angle {
	return s1.Angle(km / EarthRadiusKm)
}

// AngleToKm ...
func AngleToKm(angle s1.Angle) float64 {
	return EarthRadiusKm * float64(angle)
}

// SplitterOpts ...
type SplitterOpts struct {
	Delimiter string
	MinVars   int
	MaxVars   int
}

// PointSplitterOpts ...
var PointSplitterOpts = &SplitterOpts{
	MinVars:   2,
	MaxVars:   3,
	Delimiter: PointDelimiter,
}

// Splitter ...
func Splitter(s string, opts *SplitterOpts, vars ...*string) error {
	parts := strings.Split(s, opts.Delimiter)
	switch {
	case len(vars) < opts.MinVars:
		return errors.New("Too few variable pointers were passed")
	case len(vars) > opts.MaxVars:
		return errors.New("Too many variable pointers were passed")
	case len(parts) > opts.MaxVars:
		return errors.New("Too many delimited items")
	default:
		for i, str := range parts {
			*vars[i] = str
		}
	}
	return nil
}

// ConvertLatLon ...
func ConvertLatLon(lat, lon string) (float64, float64) {
	lat64, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Fatal(err)
	}
	lon64, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Fatal(err)
	}
	return lat64, lon64
}
