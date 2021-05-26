package s2util

import (
	"bytes"
	"encoding/json"

	"github.com/golang/geo/s2"
)

type MultiPoint []Point

func (MultiPoint) Type() string {
	return "MultiPoint"
}

func (cds MultiPoint) Point() Point {
	if len(cds) == 0 {
		panic(cds)
	}
	return cds[0]
}

func (cds MultiPoint) S2Point() s2.Point {
	return cds.Point().S2Point()
}

func (cds MultiPoint) S2Region() s2.Region {
	return nil
}

func (cds *MultiPoint) UnmarshalText(str []byte) error {
	for _, s := range bytes.Split(str, []byte(`/`)) {
		if len(s) != 0 {
			var p Point
			err := p.UnmarshalText(s)
			if err != nil {
				return err
			}
			*cds = append(*cds, p)
		}
	}
	return nil
}

func (cds *MultiPoint) UnmarshalJSON(str []byte) error {
	var v []Point

	if err := json.Unmarshal(str, &v); err != nil {
		return err
	}

	*cds = v
	return nil
}

func (cds MultiPoint) Reverse() MultiPoint {
	for i, j := 0, len(cds)-1; i < j; i, j = i+1, j-1 {
		cds[i], cds[j] = cds[j], cds[i]
	}
	return cds
}

func (cds MultiPoint) Radiusp() *float64 {
	return nil
}
