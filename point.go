package s2util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

type Point struct {
	lat Angle
	lng Angle
	alt *float64 // altitude
}

func (Point) Type() string {
	return "Point"
}

func NewPoint(lat, lng Angle, altitude *float64) (latlongalt Point) {
	latlongalt.lat = lat
	latlongalt.lng = lng
	latlongalt.alt = altitude
	return
}

func (latlong *Point) UnmarshalText(iso6709 []byte) error {
	re := regexp.MustCompile(`(?P<Latitude>[\+-][\d.]+)(?P<Longitude>[\+-][\d.]+)(?P<Altitude>[\+-][\d.]+)?`)

	if re.Match(iso6709) {
		match := re.FindSubmatch(iso6709)

		var lat, lng Angle
		var altitude *float64

		for i, name := range re.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			switch name {
			case "Latitude":
				lat = AngleFromBytes(match[i])
			case "Longitude":
				lng = AngleFromBytes(match[i])
			case "Altitude":
				altitude = getAlt(match[i])
			}
		}
		*latlong = NewPoint(lat, lng, altitude)
		return nil
	}
	panic(iso6709)
}

func NewPointFromS2Point(p s2.Point) Point {
	s2ll := s2.LatLngFromPoint(p)
	return Point{lat: NewAngleFromS1Angle(s2ll.Lat, 0), lng: NewAngleFromS1Angle(s2ll.Lng, 0)}
}

func (latlong Point) Lat() Angle {
	return latlong.lat
}

func (latlong Point) Lng() Angle {
	return latlong.lng
}

func (latlong Point) S2LatLng() s2.LatLng {
	return s2.LatLng{Lat: latlong.Lat().S1Angle(), Lng: latlong.Lng().S1Angle()}
}

func (latlong Point) S2Point() s2.Point {
	return s2.PointFromLatLng(latlong.S2LatLng())
}

func (latlong Point) S2Region() s2.Region {
	return s2.PointFromLatLng(latlong.S2LatLng())
}

func (latlong Point) Radiusp() *float64 {
	return nil
}

func (latlong Point) DistanceAngle(latlong1 *Point) s1.Angle {
	return latlong.S2LatLng().Distance(latlong1.S2LatLng())
}

func (latlong Point) DistanceEarthKm(latlong1 *Point) Km {
	return EarthArcFromAngle(latlong.DistanceAngle(latlong1))
}

func (latlong Point) latString() string {
	return strconv.FormatFloat(latlong.Lat().Degrees(), 'f', latlong.lat.preclog(), 64)
}

func (latlong Point) lngString() string {
	return strconv.FormatFloat(latlong.Lng().Degrees(), 'f', latlong.lng.preclog(), 64)
}

func getAlt(part []byte) (altitude *float64) {
	part = bytes.TrimSpace(part)
	if a, er := strconv.ParseFloat(string(part), 64); er == nil {
		altitude = &a
	}
	return
}

func (latlong Point) PrecisionArea() float64 {
	return latlong.lat.PrecDegrees() * latlong.lng.PrecDegrees()
}

func (latlong Point) PrecString() (s string) {
	s = fmt.Sprintf("lat. error %fdeg., long. error %fdeg.", latlong.lat.PrecDegrees(), latlong.lng.PrecDegrees())
	return
}

func (latlong Point) MarshalJSON() ([]byte, error) {
	var ll []Angle

	if latlong.alt != nil {
		ll = make([]Angle, 3)
		ll[2].radian = s1.Angle(*latlong.alt) * s1.Degree
		ll[2].radianprec = 1
	} else {
		ll = make([]Angle, 2)
	}

	ll[0] = latlong.lng
	ll[1] = latlong.lat

	return json.Marshal(&ll)
}

func (latlong *Point) UnmarshalJSON(data []byte) (err error) {
	var ll []Angle

	err = json.Unmarshal(bytes.TrimSpace(data), &ll)

	if len(ll) < 2 {
		return errors.New("unknown JSON Coordinate format")
	}

	latlong.lng = ll[0]
	latlong.lat = ll[1]

	if len(ll) > 2 {
		altitude := ll[2].radian.Degrees()
		latlong.alt = &altitude
	}

	return
}

func (latlong Point) altString() string {
	if latlong.alt != nil {
		return strconv.FormatFloat(*latlong.alt, 'f', 0, 64)
	}
	return ""
}
