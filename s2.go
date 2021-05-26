package s2util

import (
	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

func EdgeIntersection(a, b s2.Edge) s2.Point {
	va := s2.Point{Vector: a.V0.PointCross(a.V1).Normalize()}
	vb := s2.Point{Vector: b.V0.PointCross(b.V1).Normalize()}
	x := s2.Point{Vector: va.PointCross(vb).Normalize()}

	if v1, v2 := a.V0.Add(a.V1.Vector), b.V0.Add(b.V1.Vector); x.Dot(v1.Add(v2)) < 0 {
		x = s2.Point{Vector: r3.Vector{X: -x.X, Y: -x.Y, Z: -x.Z}}
	}

	if s2.OrderedCCW(a.V0, x, a.V1, va) && s2.OrderedCCW(b.V0, x, b.V1, vb) {
		return x
	}

	dmin2, vmin := 10.0, x
	findMinDist := func(y s2.Point) {
		d2 := x.Sub(y.Vector).Norm2()
		if d2 < dmin2 || (d2 == dmin2 && y.Cmp(vmin.Vector) == -1) {
			dmin2, vmin = d2, y
		}
	}
	if s2.OrderedCCW(b.V0, a.V0, b.V1, vb) {
		findMinDist(a.V0)
	}
	if s2.OrderedCCW(b.V0, a.V1, b.V1, vb) {
		findMinDist(a.V1)
	}
	if s2.OrderedCCW(a.V0, b.V0, a.V1, va) {
		findMinDist(b.V0)
	}
	if s2.OrderedCCW(a.V0, b.V1, a.V1, va) {
		findMinDist(b.V1)
	}
	return vmin
}

type circularLoop struct {
	s2.Point
	Intersection bool
	Done         bool

	next, prev *circularLoop
}

func newCircularLoop(p s2.Point) *circularLoop {
	c := &circularLoop{Point: p}
	c.next = c
	c.prev = c
	return c
}

func circularLoopFromPoints(pts []s2.Point) *circularLoop {
	if len(pts) == 0 {
		return circularLoopFromPoints([]s2.Point{s2.OriginPoint()})
	}

	c := newCircularLoop(pts[0])
	d := c
	for _, p := range pts[1:] {
		d = d.push(p)
	}
	return c
}

func circularLoopFromCell(cell s2.Cell) *circularLoop {
	c := newCircularLoop(cell.Vertex(0))
	d := c
	for i := 1; i < 4; i++ {
		d = d.push(cell.Vertex(i))
	}
	return c
}

func (c *circularLoop) PushIntersection(p s2.Point) {
	if e := c.Next(); e.Intersection && c.Distance(e.Point) < c.Distance(p) {
		e.PushIntersection(p)
	} else {
		c.push(p).Intersection = true
	}
}

func (c *circularLoop) Find(p s2.Point) *circularLoop {
	for d := c; ; {
		if d.Point == p {
			return d
		}
		if d = d.Next(); d == c {
			break
		}
	}
	return nil
}

func (c *circularLoop) Del() *circularLoop {
	b := c.Prev()
	d := c.Next()

	c.prev = nil
	c.next = nil
	b.next = d
	d.prev = b

	return b
}

func (c *circularLoop) Next() *circularLoop { return c.next }

func (c *circularLoop) Prev() *circularLoop { return c.prev }

func (c *circularLoop) Do(fn func(*circularLoop)) {
	fn(c)
	for p := c.Next(); p != c; p = p.Next() {
		fn(p)
	}
}

func (c *circularLoop) DoEdges(fn func(*circularLoop, *circularLoop)) {
	first := c.nextVertex()
	for d := first; ; {
		e := d.Next().nextVertex()
		fn(d, e)

		if d = e; d == c {
			break
		}
	}
}

func (c *circularLoop) push(p s2.Point) *circularLoop {
	e := c.Next()
	d := &circularLoop{Point: p}
	c.next = d
	d.prev = c
	d.next = e
	e.prev = d
	return d
}

func (c *circularLoop) nextVertex() *circularLoop {
	for d := c; ; {
		if !d.Intersection {
			return d
		}

		if d = d.Next(); d == c {
			break
		}
	}
	return nil
}

func LoopIntersectionWithCell(loop *s2.Loop, cell s2.Cell) []*s2.Loop {

	if wrap := s2.LoopFromCell(cell); loop.ContainsCell(cell) {
		return []*s2.Loop{wrap}
	} else if wrap.Contains(loop) {
		return []*s2.Loop{loop}
	}

	if !loop.IntersectsCell(cell) {
		return nil
	}

	subj, clip := circularLoopFromCell(cell), circularLoopFromPoints(loop.Vertices())

	subj.DoEdges(func(a0, a1 *circularLoop) {
		crosser := s2.NewEdgeCrosser(a0.Point, a1.Point)

		clip.DoEdges(func(b0, b1 *circularLoop) {
			if crosser.EdgeOrVertexCrossing(b0.Point, b1.Point) {
				x := EdgeIntersection(s2.Edge{V0: a0.Point, V1: a1.Point}, s2.Edge{V0: b0.Point, V1: b1.Point})
				a0.PushIntersection(x)
				b0.PushIntersection(x)
			}
		})
	})

	var res []*s2.Loop

	clip.Do(func(p *circularLoop) {
		if !p.Done && p.Intersection && !cell.ContainsPoint(p.Prev().Point) && cell.ContainsPoint(p.Next().Point) {
			pts := make([]s2.Point, 0, 4)

			c1, c2 := p, subj
			for i := 0; ; i++ {
				pts = append(pts, c1.Point)
				c1.Done = true

				if c1 = c1.Next(); c1.Intersection {
					c1, c2 = c2.Find(c1.Point), c1
				}
				if c1 == p {
					break
				}
			}

			res = append(res, s2.LoopFromPoints(pts))
		}
	})

	return res
}

func FitLoop(loop *s2.Loop, acc s2.CellUnion, maxLevel int) s2.CellUnion {
	FitLoopDo(loop, maxLevel, func(cellID s2.CellID) bool {
		acc = append(acc, cellID)
		return true
	})
	return acc
}

func FitLoopDo(loop *s2.Loop, maxLevel int, fn func(s2.CellID) bool) {
	for i := 0; i < 6; i++ {
		cellID := s2.CellIDFromFace(i)
		if nxt := fitLoopDo(loop, cellID, maxLevel, fn); !nxt {
			return
		}
	}
}

func fitLoopDo(loop *s2.Loop, cellID s2.CellID, maxLevel int, fn func(s2.CellID) bool) bool {
	cell := s2.CellFromCellID(cellID)

	if loop.ContainsCell(cell) {
		return fn(cellID)
	} else if loop.IntersectsCell(cell) {
		if cell.Level() == maxLevel {
			return fn(cellID)
		} else {
			for _, childID := range cellID.Children() {
				if !fitLoopDo(loop, childID, maxLevel, fn) {
					return false
				}
			}
		}
	}
	return true
}
