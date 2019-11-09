// Package tool copied from https://gist.github.com/antoniomo/3371e44cbe2f0cc75a525aac0d188cfb
package tool

import (
	"strconv"
	"strings"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
)

// https://blog.nobugware.com/post/2016/geo_db_s2_geohash_database/
// http://s2geometry.io/devguide/cpp/quickstart

const (
	// The Earth's mean radius in kilometers (according to NASA).
	earthRadiusKm  = 6371.01
	pointDelimiter = "::"
)

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

// PointsInCellID ...
func PointsInCellID(s2cap s2.Cap, cov s2.CellID, center Point, points []Point) {
	bmin := uint64(cov.RangeMin())
	bmax := uint64(cov.RangeMax())

	for _, v := range points {
		// This simulates an indexed range query on the DB
		if uint64(v.CellID) < bmin || uint64(v.CellID) > bmax {
			continue
		}
		// Only those in range
		ll := v.CellID.LatLng()
		lat := ll.Lat.Degrees()
		lon := ll.Lng.Degrees()
		log.Infof("Nearby Candidate: %f (lat) %f (lon) %s", lat, lon, v.Name)
		log.Infof("Calculated distance to %s: %f (km)", center.Name,
			AngleToKm(ll.Distance(center.LatLon)))
		log.Info("False positive? ", !s2cap.ContainsPoint(v.CellID.Point()))
	}
}

// KmToAngle converts a distance on the Earth's surface to an angle.
// https://github.com/golang/geo/blob/23949e136d58aeb8aa39844a312b68d90c4eb8aa/s2/s2_test.go#L38-L43
func KmToAngle(km float64) s1.Angle {
	return s1.Angle(km / earthRadiusKm)
}

// AngleToKm ...
func AngleToKm(angle s1.Angle) float64 {
	return earthRadiusKm * float64(angle)
}

// Vars 'n' stuff
var (
	// https://www.movable-type.co.uk/scripts/latlong.html
	Points = []Point{
		NewPoint(60.2934, 25.0378, "Vantaa Center"),
		NewPoint(60.2055, 24.6559, "Espoo Center"),
		NewPoint(60.1699, 24.9380, "Person in Helsinki"),
		NewPoint(50.0, 150.0, "far"),
		NewPoint(150.0, 50.0, "far"),
		NewPoint(150.0, 150.0, "far"),
		NewPoint(50.0, -50.0, "far"),
	}
)

// Run https://godoc.org/github.com/golang/geo/s2#Cap
func Run(center string) {
	c := ParsePoint(center)
	log.Infof("Center cell id: %#v", c)
	Points = append(Points, c)

	s2cap := s2.CapFromCenterAngle(c.CellID.Point(), KmToAngle(12.5))
	// http://s2geometry.io/resources/s2cell_statistics.html
	// Level 12 are about 3 to 6.4km^2 cells
	// Level 20 are about 46.41 to 97.3 meter cells
	// Since we put a MaxCells of 8, it won't go to the max level if it
	// can't approximate the region better anyway.
	rc := &s2.RegionCoverer{MaxLevel: 20, MaxCells: 8}
	covering := rc.Covering(s2.Region(s2cap))

	for i, cov := range covering {
		log.Infof("Covering Cell %d ID: %d Level: %d", i, uint64(cov),
			cov.Level())
		PointsInCellID(s2cap, cov, c, Points)
	}
}

// ParsePoint ...
func ParsePoint(delimited string) Point {
	log.Debugf("Pasring '%s' ...", delimited)
	args := strings.Split(delimited, pointDelimiter)
	var err error
	var latArg, lonArg, name string
	var lat, lon float64
	switch {
	case len(args) < 2:
		log.Errorf("Bad delimited point: %s", delimited)
		log.Fatal("At a minium, latitude and longitude must be provided")
	case len(args) == 2:
		latArg = args[0]
		lonArg = args[1]
		name = ""
	default:
		latArg = args[0]
		lonArg = args[1]
		name = args[2]
	}
	lat, err = strconv.ParseFloat(latArg, 64)
	if err != nil {
		log.Fatal(err)
	}
	lon, err = strconv.ParseFloat(lonArg, 64)
	if err != nil {
		log.Fatal(err)
	}
	return NewPoint(lat, lon, name)
}
