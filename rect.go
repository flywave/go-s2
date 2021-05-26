package s2util

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"unicode"

	geohash "github.com/TomiHiltunen/geohash-golang"
	"github.com/golang/geo/s2"
)

type Rect struct {
	s2.Rect
}

func (rect *Rect) MarshalJSON() (bb []byte, e error) {
	type LatLngs []Point

	v := []Point{
		{
			lat: NewAngleFromS1Angle(rect.Rect.Vertex(0).Lat, rect.Rect.Size().Lat/10),
			lng: NewAngleFromS1Angle(rect.Rect.Vertex(0).Lng, rect.Rect.Size().Lng/10)},
		{
			lat: NewAngleFromS1Angle(rect.Rect.Vertex(1).Lat, rect.Rect.Size().Lat/10),
			lng: NewAngleFromS1Angle(rect.Rect.Vertex(1).Lng, rect.Rect.Size().Lng/10)},
		{
			lat: NewAngleFromS1Angle(rect.Rect.Vertex(2).Lat, rect.Rect.Size().Lat/10),
			lng: NewAngleFromS1Angle(rect.Rect.Vertex(2).Lng, rect.Rect.Size().Lng/10)},
		{
			lat: NewAngleFromS1Angle(rect.Rect.Vertex(3).Lat, rect.Rect.Size().Lat/10),
			lng: NewAngleFromS1Angle(rect.Rect.Vertex(3).Lng, rect.Rect.Size().Lng/10)},
	}

	bs := make([][]byte, 0)

	for i := range v {
		b, e := v[i].MarshalJSON()
		if e != nil {
			break
		}
		bs = append(bs, b)
	}

	bb = append(bb, '[')
	bb = append(bb, bytes.Join(bs, []byte(","))...)
	bb = append(bb, ']')
	return bb, e
}

func NewRect(latitude, longitude, latprec, longprec float64) *Rect {
	rect := new(Rect)
	rect.Rect = s2.RectFromCenterSize(
		s2.LatLngFromDegrees(latitude, longitude),
		s2.LatLngFromDegrees(latprec, longprec))
	return rect
}

func NewRectGridLocator(gl string) *Rect {
	latitude := float64(-90)
	longitude := float64(-180)

	latprec := float64(10) * 24
	lonprec := float64(20) * 24

loop:
	for i, c := range gl {
		c = unicode.ToUpper(c)
		switch i % 4 {
		case 0:
			if unicode.IsUpper(c) {
				lonprec /= 24
				longitude += lonprec * float64(c-'A')
			} else {
				break loop
			}
		case 1:
			if unicode.IsUpper(c) {
				latprec /= 24
				latitude += latprec * float64(c-'A')
			} else {
				break loop
			}
		case 2:
			if unicode.IsDigit(c) {
				lonprec /= 10
				longitude += lonprec * float64(c-'0')
			} else {
				break loop
			}
		case 3:
			if unicode.IsDigit(c) {
				latprec /= 10
				latitude += latprec * float64(c-'0')
			} else {
				break loop
			}
		}
	}
	return NewRect(latitude+latprec/2, longitude+lonprec/2, latprec, lonprec)
}

func (rect Rect) Center() *Point {
	return &Point{
		lat: NewAngleFromS1Angle(rect.Rect.Center().Lat, rect.Rect.Size().Lat/2),
		lng: NewAngleFromS1Angle(rect.Rect.Center().Lng, rect.Rect.Size().Lng/2)}
}

func (rect Rect) PrecString() (s string) {
	s = fmt.Sprintf("lat. error %fdeg., long. error %fdeg.", rect.Size().Lat.Degrees(), rect.Size().Lng.Degrees())
	return
}

func (rect *Rect) GridLocator() string {
	const floaterr = 1 + 1e-11

	var gl []rune

	latitude := rect.Center().Lat().Degrees() + 90
	longitude := rect.Center().Lng().Degrees() + 180

	latprec := float64(10) * 24
	lonprec := float64(20) * 24

loop:
	for i := 0; ; i++ {
		switch i % 4 {
		case 0:
			lonprec /= 24
			if lonprec*floaterr < rect.Size().Lng.Degrees() {
				break loop
			}
			c := math.Floor(longitude / lonprec)
			gl = append(gl, rune(byte(c)+'A'))
			longitude -= c * lonprec
		case 1:
			latprec /= 24
			if latprec*floaterr < rect.Size().Lat.Degrees() {
				break loop
			}
			c := math.Floor(latitude / latprec)
			gl = append(gl, rune(byte(c)+'A'))
			latitude -= c * latprec
		case 2:
			lonprec /= 10
			if lonprec*floaterr < rect.Size().Lng.Degrees() {
				break loop
			}
			c := math.Floor(longitude / lonprec)
			gl = append(gl, rune(byte(c)+'0'))
			longitude -= c * lonprec

		case 3:
			latprec /= 10
			if latprec*floaterr < rect.Size().Lat.Degrees() {
				break loop
			}
			c := math.Floor(latitude / latprec)
			gl = append(gl, rune(byte(c)+'0'))
			latitude -= c * latprec
		}
	}

	l := len(gl)
	if l%2 == 1 {
		gl = gl[:l-1]
	}
	return string(gl)
}

func NewRectGeoHash(geoHash string) (latlong *Rect, err error) {
	if bb := geohash.Decode(geoHash); bb != nil {
		latlong = NewRect(bb.Center().Lat(), bb.Center().Lng(), bb.NorthEast().Lat()-bb.SouthWest().Lat(), bb.NorthEast().Lng()-bb.SouthWest().Lng())
	} else {
		err = errors.New("Geohash decode error")
	}
	return
}

func (rect *Rect) geoHash(precision int) string {
	return geohash.EncodeWithPrecision(rect.Center().Lat().Degrees(), rect.Center().Lng().Degrees(), precision)
}

func (rect *Rect) GeoHash5() string {
	return rect.geoHash(5)
}

func (rect *Rect) GeoHash6() string {
	return rect.geoHash(6)
}

func (rect *Rect) GeoHash() string {
	const floaterr = 1 + 5e-10

	geohashlatbits := -math.Log2(rect.Size().Lat.Degrees()/45) + 2 // div by 180 = 45 * 2^2
	geohashlngbits := -math.Log2(rect.Size().Lng.Degrees()/45) + 3 // div by 360 = 45 * 2^3
	geohashlat2len, geohashlatlen2mod := math.Modf(geohashlatbits / 5 * floaterr)

	var geohashlatlen int
	if geohashlatlen2mod >= 0.4 {
		geohashlatlen = int(geohashlat2len)*2 + 1
	} else {
		geohashlatlen = int(geohashlat2len) * 2
	}

	geohashlng2len, geohashlnglen2mod := math.Modf(geohashlngbits / 5 * floaterr)

	var geohashlnglen int
	if geohashlnglen2mod >= 0.6 {
		geohashlnglen = int(geohashlng2len)*2 + 1
	} else {
		geohashlnglen = int(geohashlng2len) * 2
	}

	if geohashlatlen < geohashlnglen {
		return rect.geoHash(geohashlatlen)
	}
	return rect.geoHash(geohashlnglen)
}

func (rect *Rect) S2Rect() s2.Rect {
	return rect.Rect
}

func (rect *Rect) S2Region() s2.Region {
	return rect.S2Rect()
}
