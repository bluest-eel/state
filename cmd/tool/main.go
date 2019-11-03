// Copied from https://gist.github.com/antoniomo/3371e44cbe2f0cc75a525aac0d188cfb
package main

import (
	"fmt"

	"github.com/bluest-eel/state/tool"
	"github.com/golang/geo/s2"
)

// Vars 'n' stuff
var (
	LLH = s2.LatLngFromDegrees(60.1699, 24.9384) // Helsinki Center
	// https://www.movable-type.co.uk/scripts/latlong.html
	Points = []tool.Point{
		tool.NewPoint(60.1699, 24.9384, "Helsinki Center"),
		tool.NewPoint(60.2934, 25.0378, "Vantaa Center (14.79km)"),
		tool.NewPoint(60.2055, 24.6559, "Espoo Center (16.11km)"),
		tool.NewPoint(60.1699, 24.9380, "Person in Helsinki (22m)"),
		tool.NewPoint(50.0, 150.0, "far"),
		tool.NewPoint(150.0, 50.0, "far"),
		tool.NewPoint(150.0, 150.0, "far"),
		tool.NewPoint(50.0, -50.0, "far"),
	}
)

// https://godoc.org/github.com/golang/geo/s2#Cap
func main() {
	c := s2.CellIDFromLatLng(LLH)
	fmt.Println(c)

	s2cap := s2.CapFromCenterAngle(c.Point(), tool.KmToAngle(12.5))
	// http://s2geometry.io/resources/s2cell_statistics.html
	// Level 12 are about 3 to 6.4km^2 cells
	// Level 20 are about 46.41 to 97.3 meter cells
	// Since we put a MaxCells of 8, it won't go to the max level if it
	// can't approximate the region better anyway.
	rc := &s2.RegionCoverer{MaxLevel: 20, MaxCells: 8}
	covering := rc.Covering(s2.Region(s2cap))

	for i, cov := range covering {
		fmt.Printf("Covering Cell %d ID: %d Level: %d\n", i, uint64(cov), cov.Level())
		tool.PointsInCellID(s2cap, cov, LLH, Points)
	}
}
