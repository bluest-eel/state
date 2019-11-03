// Package tool copied from https://gist.github.com/antoniomo/3371e44cbe2f0cc75a525aac0d188cfb
package tool

import (
	"fmt"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
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
		fmt.Println("Nearby Candidate:", lat, lon, v.name)
		fmt.Println("Calculated distance to Helsinki Center:", AngleToKm(ll.Distance(center)))
		fmt.Println("False positive?", !s2cap.ContainsPoint(v.cellID.Point()))
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
