package s2util

import (
	"strconv"

	"github.com/golang/geo/s1"
)

const (
	radiusKmOfTheEarth = 6371.01 // radius km of the Earth.
)

type Km float64

func (km Km) EarthAngle() s1.Angle {
	return s1.Angle(km / radiusKmOfTheEarth)
}

func EarthArcFromAngle(angle s1.Angle) Km {
	return Km(angle * radiusKmOfTheEarth)
}

func (km Km) EarthChordAngle() s1.ChordAngle {
	return s1.ChordAngleFromAngle(km.EarthAngle())
}

func EarthArcFromChordAngle(chordangle s1.ChordAngle) Km {
	return EarthArcFromAngle(chordangle.Angle())
}

func (km Km) String() string {
	return strconv.FormatFloat(float64(km), 'f', 1, 64) + "km"
}
