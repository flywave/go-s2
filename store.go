package s2util

import (
	"errors"
	"io"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/golang/geo/s2"
)

var (
	errReleased = errors.New("s2store: already released")
)

type Store struct {
	db           *leveldb.DB
	storageLevel int
}

func NewStore(r io.ReaderAt, storageLevel int) (*Store, error) {
	return nil, nil
}

func (r *Store) AddLatlng(point s2.LatLng) error {
	cell := s2.CellIDFromLatLng(point)
	_ = []byte(cell.ToToken()) // token
	return nil
}

func (r *Store) AddPolygon(vertices []s2.LatLng) error {
	points := func() (res []s2.Point) {
		for _, vertex := range vertices {
			point := s2.PointFromLatLng(vertex)
			res = append(res, point)
		}
		return
	}()

	loop := s2.LoopFromPoints(points)
	loop.Normalize()
	_ = s2.PolygonFromLoops([]*s2.Loop{loop}) // polygon

	coverer := s2.RegionCoverer{MinLevel: r.storageLevel, MaxLevel: r.storageLevel}
	_ = coverer.Covering(loop) // cells
	return nil
}

func (r *Store) FindSection(cellID s2.CellID) (iterator.Iterator, error) {
	return nil, nil
}

func (r *Store) Match(p0, p1 s2.LatLng) bool {
	boundRect := s2.RectFromLatLng(p0)
	boundRect = boundRect.AddPoint(p1)

	_ = []byte(s2.CellIDFromLatLng(boundRect.Hi()).ToToken()) //start
	_ = []byte(s2.CellIDFromLatLng(boundRect.Lo()).ToToken()) // end

	return true
}

func (r *Store) Search(p0, p1 s2.LatLng) []s2.CellID {
	boundRect := s2.RectFromLatLng(p0)
	boundRect = boundRect.AddPoint(p1)

	_ = []byte(s2.CellIDFromLatLng(boundRect.Hi()).ToToken()) //start
	_ = []byte(s2.CellIDFromLatLng(boundRect.Lo()).ToToken()) // end

	results := []s2.CellID{}

	return results
}

func (r *Store) Nearby(cellID s2.CellID, limit int) (*NearbyRS, error) {
	return nil, nil
}
