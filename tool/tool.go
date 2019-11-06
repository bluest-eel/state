// Package tool copied from https://gist.github.com/antoniomo/3371e44cbe2f0cc75a525aac0d188cfb
package tool

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
)

// https://blog.nobugware.com/post/2016/geo_db_s2_geohash_database/
// http://s2geometry.io/devguide/cpp/quickstart

const (
	// The Earth's mean radius in kilometers (according to NASA).
	earthRadiusKm = 6371.01
)

// Point ...
type Point struct {
	cellID s2.CellID
	name   string
}

// NewPoint ...
func NewPoint(lat, lon float64, name string) Point {
	return Point{
		cellID: s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lon)),
		name:   name,
	}
}

// PointsInCellID ...
func PointsInCellID(s2cap s2.Cap, cov s2.CellID, center s2.LatLng, points []Point) {
	bmin := uint64(cov.RangeMin())
	bmax := uint64(cov.RangeMax())

	for _, v := range points {
		// This simulates an indexed range query on the DB
		if uint64(v.cellID) < bmin || uint64(v.cellID) > bmax {
			continue
		}
		// Only those in range
		ll := v.cellID.LatLng()
		lat := ll.Lat.Degrees()
		lon := ll.Lng.Degrees()
		log.Infof("Nearby Candidate: %f (lat) %f (lon) %s", lat, lon, v.name)
		log.Infof("Calculated distance to Helsinki Center: %f km",
			AngleToKm(ll.Distance(center)))
		log.Info("False positive? ", !s2cap.ContainsPoint(v.cellID.Point()))
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
	LLH = s2.LatLngFromDegrees(60.1699, 24.9384) // Helsinki Center
	// https://www.movable-type.co.uk/scripts/latlong.html
	Points = []Point{
		NewPoint(60.1699, 24.9384, "Helsinki Center (0 km)"),
		NewPoint(60.2934, 25.0378, "Vantaa Center (14.79 km)"),
		NewPoint(60.2055, 24.6559, "Espoo Center (16.11 km)"),
		NewPoint(60.1699, 24.9380, "Person in Helsinki (22 m)"),
		NewPoint(50.0, 150.0, "far"),
		NewPoint(150.0, 50.0, "far"),
		NewPoint(150.0, 150.0, "far"),
		NewPoint(50.0, -50.0, "far"),
	}
)

// Run https://godoc.org/github.com/golang/geo/s2#Cap
func Run() {
	c := s2.CellIDFromLatLng(LLH)
	log.Infof("Center cell id: %#v", c)

	s2cap := s2.CapFromCenterAngle(c.Point(), KmToAngle(12.5))
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
		PointsInCellID(s2cap, cov, LLH, Points)
	}
}
