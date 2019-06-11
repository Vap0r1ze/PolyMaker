package shapes

import (
	"strconv"

	"github.com/fogleman/delaunay"
)

type Point struct {
	X int
	Y int
}

type Triangle struct {
	P1 Point
	P2 Point
	P3 Point
}

func FromTriangulation(t *delaunay.Triangulation) []Triangle {
	ts := make([]Triangle, len(t.Triangles)/3)

	index := 0
	for i := 0; i < len(t.Triangles); i += 3 {
		P1 := t.Points[t.Triangles[i+0]]
		P2 := t.Points[t.Triangles[i+1]]
		P3 := t.Points[t.Triangles[i+2]]
		ts[index] = Triangle{
			P1: Point{
				X: int(P1.X),
				Y: int(P1.Y),
			},
			P2: Point{
				X: int(P2.X),
				Y: int(P2.Y),
			},
			P3: Point{
				X: int(P3.X),
				Y: int(P3.Y),
			},
		}
		index++
	}

	return ts
}

func (t *Triangle) ContainsPoint(p Point) bool {
	s := t.P1.Y*t.P3.X - t.P1.X*t.P3.Y + (t.P3.Y-t.P1.Y)*p.X + (t.P1.X-t.P3.X)*p.Y
	z := t.P1.X*t.P2.Y - t.P1.Y*t.P2.X + (t.P1.Y-t.P2.Y)*p.X + (t.P2.X-t.P1.X)*p.Y

	if (s < 0) != (z < 0) {
		return false
	}
	A := -t.P2.Y*t.P3.X + t.P1.Y*(t.P3.X-t.P2.X) + t.P1.X*(t.P2.Y-t.P3.Y) + t.P2.X*t.P3.Y

	var f bool
	if A < 0 {
		f = (s <= 0) && (s+z >= A)
	} else {
		f = (s >= 0) && (s+z <= A)
	}

	return f
}

func (t *Triangle) GetPathData() string {
	var d string
	d += "M" + strconv.FormatInt(int64(t.P1.X), 10) + "," + strconv.FormatInt(int64(t.P1.Y), 10)
	d += "L" + strconv.FormatInt(int64(t.P2.X), 10) + "," + strconv.FormatInt(int64(t.P2.Y), 10)
	d += "L" + strconv.FormatInt(int64(t.P3.X), 10) + "," + strconv.FormatInt(int64(t.P3.Y), 10)
	d += "Z"
	return d
}
