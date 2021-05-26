package s2util

// #cgo LDFLAGS: -lm
// #include "geodesic.h"
import "C"

func InvertAzimuth(azDeg float64) float64 {
	return azDeg + 180%360
}

type Geodesic struct {
	g C.struct_geod_geodesic
}

func NewGeodesic(a float64, f float64) Geodesic {
	var g Geodesic
	C.geod_init(&g.g, C.double(a), C.double(f))
	return g
}

func (g Geodesic) Direct(lat, lon float64, az1 float64, s12 float64) (lat2, lon2, az2 float64) {

	var resLat, resLon, resAz C.double

	C.geod_direct(
		&g.g,
		C.double(lat),
		C.double(lon),
		C.double(az1),
		C.double(s12),
		&resLat,
		&resLon,
		&resAz)

    lat2 = float64(resLat)
	lon2 = float64(resLon)
	az2 = float64(resAz)

	return
}

func (g *Geodesic) Inverse(lat1, lon1, lat2, lon2 float64) (s12, az1, az2 float64) {

	var resS12, resAz1, resAz2 C.double

	C.geod_inverse(
		&g.g,
		C.double(lat1),
		C.double(lon1),
		C.double(lat2),
		C.double(lon2),
		&resS12,
		&resAz1,
		&resAz2)

	s12 = float64(resS12)
	az1 = float64(resAz1)
	az2 = float64(resAz2)

	return
}

type GeodesicLine struct {
	gl C.struct_geod_geodesicline
}

func NewGeodesicLine(g Geodesic, lat, lon float64, azDeg float64) GeodesicLine {
	var l GeodesicLine
	C.geod_lineinit(&l.gl, &g.g, C.double(lat), C.double(lon), C.double(azDeg), 0)
	return l
}

func (l GeodesicLine) Position(s12 float64) (lat, lon, azDeg float64) {

	var resLat, resLon, resAzDeg C.double

	C.geod_position(
		&l.gl,
		C.double(s12),
		&resLat,
		&resLon,
		&resAzDeg)

	lat = float64(resLat)
	lon = float64(resLon)
	azDeg = float64(resAzDeg)

	return
}
