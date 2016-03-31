package mesh

import (
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

type Box struct {
	LowerBound, UpperBound vec3.T
}

func (this Box) BoundTriangle(tri Triangle) {
	this.LowerBound = vec3.T{
		math.Min(tri[0][0], math.Min(tri[0][1], tri[0][2])),
		math.Min(tri[1][0], math.Min(tri[1][1], tri[1][2])),
		math.Min(tri[2][0], math.Min(tri[2][1], tri[2][2])),
	}

	this.UpperBound = vec3.T{
		math.Max(tri[0][0], math.Max(tri[0][1], tri[0][2])),
		math.Max(tri[1][0], math.Max(tri[1][1], tri[1][2])),
		math.Max(tri[2][0], math.Max(tri[2][1], tri[2][2])),
	}
}

func (this Box) BoundTriangles(tris ...Triangle) {
	if len(tris) == 0 {
		return
	}

	this.BoundTriangle(tris[0])
	for _, tri := range tris[1:] {
		for vert := 0; vert < 3; vert++ {
			for compon := 0; compon < 3; compon++ {
				val := tri[vert][compon]
				if val < this.LowerBound[compon] {
					this.LowerBound[compon] = val
				} else if val > this.UpperBound[compon] {
					this.UpperBound[compon] = val
				}
			}
		}
	}
}
