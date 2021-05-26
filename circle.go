package s2util

import (
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

type Circle struct {
	Point
	s1.ChordAngle
}

func (Circle) Type() string {
	return "Circle"
}

func NewCircle(latlng Point, km Km) *Circle {
	circle := Circle{
		Point:      latlng,
		ChordAngle: s1.ChordAngleFromAngle(km.EarthAngle()),
	}
	return &circle
}

func NewPointCircle(latlng Point) *Circle {
	latprecchordangle := s1.ChordAngleFromAngle(latlng.lat.PrecS1Angle())
	lngprecchordangle := s1.ChordAngleFromAngle(latlng.lng.PrecS1Angle()) * s1.ChordAngle(math.Abs(math.Cos(float64(latlng.lat.PrecS1Angle()))))
	var precchordangle s1.ChordAngle

	if latprecchordangle > lngprecchordangle {
		precchordangle = latprecchordangle
	} else {
		precchordangle = lngprecchordangle
	}

	circle := Circle{
		Point:      latlng,
		ChordAngle: precchordangle,
	}
	return &circle
}

func NewEmptyCircle() *Circle {
	circle := Circle{
		ChordAngle: s1.NegativeChordAngle,
	}
	return &circle
}

func (c Circle) S2Cap() s2.Cap {
	return s2.CapFromCenterChordAngle(s2.PointFromLatLng(c.Point.S2LatLng()), c.ChordAngle)
}

func (c Circle) S2Region() s2.Region {
	return c.S2Cap()
}

func (c *Circle) CapBound() s2.Cap {
	return c.S2Cap().CapBound()
}

func (c *Circle) RectBound() s2.Rect {
	return c.S2Cap().RectBound()
}

func (c *Circle) ContainsCell(cell s2.Cell) bool {
	return c.S2Cap().ContainsCell(cell)
}

func (c *Circle) IntersectsCell(cell s2.Cell) bool {
	return c.S2Cap().IntersectsCell(cell)
}

func (c *Circle) ContainsPoint(p s2.Point) bool {
	return c.S2Cap().ContainsPoint(p)
}

func (c *Circle) CellUnionBound() []s2.CellID {
	return c.S2Cap().CellUnionBound()
}

func (c Circle) Radiusp() *float64 {
	r := float64(c.Radius())
	return &r
}

func (c Circle) Radius() Km {
	return EarthArcFromChordAngle(c.ChordAngle)
}

func (c Circle) S2Point() s2.Point {
	return c.Point.S2Point()
}

func (c *Circle) S2Loop(div int) (loop *s2.Loop) {
	return s2.RegularLoop(s2.PointFromLatLng(c.Point.S2LatLng()), c.Angle(), div)
}

func (c *Circle) S2LatLngs(div int) (lls []s2.LatLng) {
	vs := c.S2Loop(div).Vertices()
	lls = make([]s2.LatLng, len(vs))
	for i := range vs {
		lls[i] = s2.LatLngFromPoint(vs[i])
	}
	return
}

func (c *Circle) MultiPoint(div int) (lls MultiPoint) {
	vs := c.S2Loop(div).Vertices()
	lls = make(MultiPoint, len(vs))
	for i := range vs {
		lls[i] = NewPointFromS2Point(vs[i])
	}
	return
}
