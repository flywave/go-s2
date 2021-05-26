package s2util

type DistUnitType int

const (
	DIST_UNIT_UNKNOW = DistUnitType(-1)
	DIST_UNIT_M      = DistUnitType(0)
	DIST_UNIT_KM     = DistUnitType(1)
	DIST_UNIT_MI     = DistUnitType(2)
	DIST_UNIT_NM     = DistUnitType(3)
	DIST_UNIT_FT     = DistUnitType(4)
)

func ParseDistUnit(s string) DistUnitType {
	if s == "m" {
		return DIST_UNIT_M
	} else if s == "km" {
		return DIST_UNIT_KM
	} else if s == "mi" {
		return DIST_UNIT_MI
	} else if s == "nm" {
		return DIST_UNIT_NM
	} else if s == "ft" {
		return DIST_UNIT_FT
	}
	return DIST_UNIT_UNKNOW
}

func unitToKM(u DistUnitType) float64 {
	switch u {
	case DIST_UNIT_M:
		return 0.001
	case DIST_UNIT_KM:
		return 1.0
	case DIST_UNIT_MI:
		return 1.609344
	case DIST_UNIT_NM:
		return 1.8520
	case DIST_UNIT_FT:
		return 0.0003048
	default:
		return 0.0
	}
}

func ConvertDistUnit(d float64, from DistUnitType, to DistUnitType) Km {
	convFactor := unitToKM(from)
	convFactor /= unitToKM(to)
	return Km(d * convFactor)
}

type EllipsoidSepc struct {
	EquatorRadius float64
	FLattening    float64
}

func (e EllipsoidSepc) PolesRadius() float64 {
	return (1.0 - e.FLattening) * e.EquatorRadius
}

var (
	WGS84_ELLIPSOID = EllipsoidSepc{6378137.0, 1.0 / 298.257223563}
	UNIT_SPHERE     = EllipsoidSepc{1.0, 0.0}
)

func GeodesicDistance(p1, p2 Point, e EllipsoidSepc) float64 {
	g := NewGeodesic(e.EquatorRadius, e.FLattening)
	dist, _, _ := g.Inverse(p1.lat.Degrees(), p1.lng.Degrees(), p2.lat.Degrees(), p2.lng.Degrees())
	return dist
}

func GeodesicPointAtDist(p Point, dist float64, azimuth float64, e EllipsoidSepc) Point {
	g := NewGeodesic(e.EquatorRadius, e.FLattening)
	lat2, lon2, _ := g.Direct(p.lat.Degrees(), p.lng.Degrees(), azimuth, dist)
	return Point{NewAngle(lat2, 0.1), NewAngle(lon2, 0.1), nil}
}
