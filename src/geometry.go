package mesh

import (
	"errors"
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

const epsilon = 1e-5

type Triangle [3]vec3.T

func (this Triangle) IntersectTriangle(other Triangle) (*Line, error) {
	planeA, planeB := this.Plane(), other.Plane()

	triApts := planeB.IntersectTriangle(this)
	triBpts := planeA.IntersectTriangle(other)
	if len(triApts) != 2 || len(triBpts) != 2 {
		return nil, ErrDontIntersect
	}

	aHasB0 := this.ContainsPoint(triBpts[0])
	aHasB1 := this.ContainsPoint(triBpts[1])
	bHasA0 := other.ContainsPoint(triApts[0])
	bHasA1 := other.ContainsPoint(triApts[1])

	var line *Line

	if aHasB0 {
		if aHasB1 {
			line = &Line{triBpts[0], triBpts[1]}
		} else if bHasA0 {
			line = &Line{triBpts[0], triApts[0]}
		} else if bHasA1 {
			line = &Line{triBpts[0], triApts[1]}
		}
	} else if aHasB1 {
		if bHasA0 {
			line = &Line{triBpts[1], triApts[0]}
		} else if bHasA1 {
			line = &Line{triBpts[1], triApts[1]}
		}
	} else if bHasA0 {
		if bHasA1 {
			line = &Line{triApts[0], triApts[1]}
		}
	}

	if line[0] == line[1] {
		return nil, ErrDontIntersect
	}

	return line, nil
}

func (this Triangle) Plane() Plane {
	v0, v1, v2 := this[0], this[1], this[2]
	normal := vec3.Cross(v1.Sub(&v0), (v2.Sub(&v0)))

	invNormal := normal.Scaled(-1)
	offset := vec3.Dot(&invNormal, &v0)

	return Plane{normal, offset}
}

var (
	ErrDontIntersect = errors.New("No intersection found")
	ErrCoplanar      = errors.New("The triangles are coplanar")
)

func IsDontIntersect(err error) bool {
	return err == ErrDontIntersect
}

func IsCoplanar(err error) bool {
	return err == ErrCoplanar
}

// http://www.blackpawn.com/texts/pointinpoly/
func (this Triangle) ContainsPoint(pt vec3.T) bool {
	return (Line{this[0], this[1]}.SameSide(this[2], pt) &&
		Line{this[0], this[2]}.SameSide(this[1], pt) &&
		Line{this[1], this[2]}.SameSide(this[0], pt))
}

type Line [2]vec3.T

func (this Line) SameSide(p0, p1 vec3.T) bool {
	lineVec := vec3.Sub(&this[1], &this[0])
	p0Vec := vec3.Sub(&p0, &this[0])
	p1Vec := vec3.Sub(&p1, &this[0])

	cp0 := vec3.Cross(&p0Vec, &lineVec)
	cp1 := vec3.Cross(&p1Vec, &lineVec)

	return vec3.Dot(&cp0, &cp1) >= 0
}

type Plane struct {
	Normal vec3.T
	Offset float64
}

func (this Plane) IntersectLine(line Line) *vec3.T {
	lineVec := vec3.Sub(&line[1], &line[0])
	invLineVec := vec3.Sub(&line[0], &line[1])

	denom := vec3.Dot(&this.Normal, &invLineVec)
	if math.Abs(denom) < epsilon {
		// Line is parallel to plane
		return nil
	}

	numer := vec3.Dot(&this.Normal, &line[0]) + this.Offset
	t := numer / denom
	if t < 0 || t > 1 {
		return nil
	}

	lineVec.Scale(t)
	intersectPt := vec3.Add(&lineVec, &line[0])
	return &intersectPt
}

func (this Plane) TriangleCrosses(tri Triangle) bool {
	sign1 := sign(vec3.Dot(&this.Normal, &tri[0]) + this.Offset)
	sign2 := sign(vec3.Dot(&this.Normal, &tri[1]) + this.Offset)
	sign3 := sign(vec3.Dot(&this.Normal, &tri[2]) + this.Offset)
	return !((sign1 == sign2 && sign2 == sign3) && sign1 != 0)
}

func sign(f float64) int {
	if math.Abs(f) < epsilon {
		return 0
	}

	if f > 0 {
		return 1
	}
	return -1
}

func (this Plane) IntersectTriangle(tri Triangle) []vec3.T {
	points := make([]vec3.T, 0, 2)
	pt0 := this.IntersectLine(Line{tri[0], tri[1]})
	pt1 := this.IntersectLine(Line{tri[1], tri[2]})
	pt2 := this.IntersectLine(Line{tri[2], tri[0]})

	if pt0 != nil {
		points = append(points, *pt0)
	}
	if pt1 != nil {
		points = append(points, *pt1)
	}
	if pt2 != nil {
		points = append(points, *pt2)
	}

	if len(points) != 3 {
		return points
	}

	if points[0] == points[1] || points[0] == points[2] {
		return points[1:]
	}
	if points[1] == points[2] {
		return points[:2]
	}

	// Shouldn't ever get here
	return points
}
