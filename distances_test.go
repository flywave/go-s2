package s2util

import "testing"

func TestDistances(t *testing.T) {
	dist := GeodesicDistance(Point{NewAngle(36., 0.0001), NewAngle(117., 0.0001), nil}, Point{NewAngle(36.8298321, 0.0001), NewAngle(117.8778809, 0.0001), nil}, WGS84_ELLIPSOID)

	if dist <= 0 {
		t.Error("error")
	}
}
