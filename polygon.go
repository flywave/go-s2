package s2util

import (
	"encoding/json"

	"github.com/golang/geo/s2"
)

type Polygon struct {
	LineString
}

func (Polygon) Type() string {
	return "Polygon"
}

func (cds Polygon) S2Loop() *s2.Loop {
	ps := make(s2.Polyline, len(cds.MultiPoint))
	for i := range cds.MultiPoint {
		ps[i] = cds.MultiPoint[i].S2Point()
	}
	l := s2.LoopFromPoints(ps)
	if !l.IsNormalized() {
		l.Invert()
	}
	return l
}

func (cds Polygon) S2Region() s2.Region {
	return cds.S2Loop()
}

func (cds *Polygon) CapBound() s2.Cap {
	return cds.S2Loop().CapBound()
}

func (cds *Polygon) RectBound() s2.Rect {
	return cds.S2Loop().RectBound()
}

func (cds *Polygon) ContainsCell(c s2.Cell) bool {
	return cds.S2Loop().ContainsCell(c)
}

func (cds *Polygon) IntersectsCell(c s2.Cell) bool {
	return cds.S2Loop().IntersectsCell(c)
}

func (cds *Polygon) ContainsPoint(p s2.Point) bool {
	return cds.S2Loop().ContainsPoint(p)
}

func (cds *Polygon) CellUnionBound() []s2.CellID {
	return cds.S2Loop().CellUnionBound()
}

func (cds Polygon) S2Point() s2.Point {
	return cds.S2Loop().Centroid()
}

func (cds Polygon) Radiusp() *float64 {
	return nil
}

func (cds Polygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[]MultiPoint{cds.MultiPoint})
}

func (cds *Polygon) UnmarshalJSON(data []byte) (err error) {
	var co []MultiPoint
	err = json.Unmarshal(data, &co)
	if err != nil {
		panic(err)
	}

	switch len(co) {
	case 0:
		panic("No Polygon!")
	case 1:
		cds.MultiPoint = co[0]
	default:
		panic("Polygon has hole! Not implemented")
	}
	return nil
}
