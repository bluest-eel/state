// Copied from https://gist.github.com/antoniomo/3371e44cbe2f0cc75a525aac0d188cfb
package main

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

type point struct {
	cellID s2.CellID
	name   string
}

func newPoint(lat, lon float64, name string) point {
	return point{
		cellID: s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lon)),
		name:   name,
	}
}

var (
	llh = s2.LatLngFromDegrees(60.1699, 24.9384) // Helsinki Center
	// https://www.movable-type.co.uk/scripts/latlong.html
	points = []point{
		newPoint(60.1699, 24.9384, "Helsinki Center"),
		newPoint(60.2934, 25.0378, "Vantaa Center (14.79km)"),
		newPoint(60.2055, 24.6559, "Espoo Center (16.11km)"),
		newPoint(60.1699, 24.9380, "Person in Helsinki (22m)"),
		newPoint(50.0, 150.0, "far"),
		newPoint(150.0, 50.0, "far"),
		newPoint(150.0, 150.0, "far"),
		newPoint(50.0, -50.0, "far"),
	}
)

func pointsInCellID(s2cap s2.Cap, cov s2.CellID, points []point) {
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
		fmt.Println("Calculated distance to Helsinki Center:", angleToKm(ll.Distance(llh)))
		fmt.Println("False positive?", !s2cap.ContainsPoint(v.cellID.Point()))
	}
}

// kmToAngle converts a distance on the Earth's surface to an angle.
// https://github.com/golang/geo/blob/23949e136d58aeb8aa39844a312b68d90c4eb8aa/s2/s2_test.go#L38-L43
func kmToAngle(km float64) s1.Angle {
	return s1.Angle(km / earthRadiusKm)
}

func angleToKm(angle s1.Angle) float64 {
	return earthRadiusKm * float64(angle)
}
