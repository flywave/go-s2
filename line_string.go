package s2util

import (
	"encoding/json"

	"github.com/golang/geo/s2"
)

type LineString struct {
	MultiPoint
}

func (LineString) Type() string {
	return "LineString"
}

func (cds *LineString) Polygon() Polygon {
	return Polygon{LineString: *cds}
}

func (cds *LineString) UnmarshalText(b []byte) error {
	return cds.MultiPoint.UnmarshalText(b)
}

func (cds LineString) S2Polyline() s2.Polyline {
	var ps s2.Polyline
	for _, cd := range cds.MultiPoint {
		ps = append(ps, cd.S2Point())
	}
	return ps
}

func (cds LineString) S2Loop() *s2.Loop {
	lo := s2.LoopFromPoints(cds.S2Polyline())
	lo.Normalize()
	return lo
}

func (cds LineString) S2Point() s2.Point {
	return cds.MultiPoint[0].S2Point()
}

func (cds LineString) S2Region() s2.Region {
	ps := make(s2.Polyline, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		ps[i] = cds.MultiPoint[i].S2Point()
	}
	return &ps
}

func (cds LineString) Radiusp() *float64 {
	return nil
}

func (cds *LineString) CapBound() s2.Cap {
	return cds.S2Region().CapBound()
}

func (cds *LineString) RectBound() s2.Rect {
	return cds.S2Region().RectBound()
}

func (cds *LineString) ContainsCell(c s2.Cell) bool {
	return cds.S2Region().ContainsCell(c)
}

func (cds *LineString) IntersectsCell(c s2.Cell) bool {
	return cds.S2Region().IntersectsCell(c)
}

func (cds *LineString) ContainsPoint(p s2.Point) bool {
	return cds.S2Region().ContainsPoint(p)
}

func (cds *LineString) CellUnionBound() []s2.CellID {
	return cds.S2Region().CellUnionBound()
}

func (cds LineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(&cds.MultiPoint)
}

func (cds *LineString) UnmarshalJSON(data []byte) (err error) {
	err = json.Unmarshal(data, &cds.MultiPoint)
	if err != nil {
		panic(err)
	}
	return nil
}
