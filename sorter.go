package s2util

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/bsm/extsort"
	"github.com/golang/geo/s2"
)

var (
	errInvalidCellID = errors.New("s2store: invalid cell ID")
)

type SorterOptions struct {
	TempDir string
}

func (o *SorterOptions) norm() *SorterOptions {
	var oo SorterOptions
	if o != nil {
		oo = *o
	}
	return &oo
}

type Sorter struct {
	x *extsort.Sorter
	t []byte
}

func NewSorter(o *SorterOptions) *Sorter {
	o = o.norm()
	return &Sorter{
		x: extsort.New(&extsort.Options{WorkDir: o.TempDir}),
	}
}

func (s *Sorter) Append(cellID s2.CellID, data []byte) error {
	if !cellID.IsValid() {
		return errInvalidCellID
	}

	if sz := 8 + len(data); sz < cap(s.t) {
		s.t = s.t[:sz]
	} else {
		s.t = make([]byte, sz)
	}

	binary.BigEndian.PutUint64(s.t[0:], uint64(cellID))
	copy(s.t[8:], data)
	return s.x.Append(s.t)
}

func (s *Sorter) Sort() (*SorterIterator, error) {
	iter, err := s.x.Sort()
	if err != nil {
		return nil, err
	}
	return &SorterIterator{it: iter}, nil
}

func (s *Sorter) Close() error {
	return s.x.Close()
}

type SorterIterator struct {
	it *extsort.Iterator

	current [][]byte
	nextID  s2.CellID
	next    [][]byte
}

func (i *SorterIterator) NextEntry() (s2.CellID, [][]byte, error) {
	currentID := i.nextID
	for i.it.Next() {
		rawdata := i.it.Data()
		i.nextID = s2.CellID(binary.BigEndian.Uint64(rawdata))

		if currentID != 0 && currentID != i.nextID {
			i.next = i.push(i.next, rawdata[8:])
			break
		}
		currentID = i.nextID
		i.current = i.push(i.current, rawdata[8:])
	}

	if err := i.it.Err(); err != nil {
		return 0, nil, err
	}

	if size := len(i.current); size != 0 {
		i.current, i.next = i.next, i.current[:0]
		return currentID, i.next[:size], nil
	}

	return 0, nil, io.EOF
}
func (i *SorterIterator) Close() error {
	return i.it.Close()
}

func (i *SorterIterator) push(chunks [][]byte, chunk []byte) [][]byte {
	if pos := len(chunks); pos+1 < cap(chunks) {
		chunks = chunks[:pos+1]
		chunks[pos] = append(chunks[pos][:0], chunk...)
	} else {
		cloned := make([]byte, len(chunk))
		copy(cloned, chunk)
		chunks = append(chunks, cloned)
	}
	return chunks
}
