package mesh

import "github.com/ungerik/go3d/vec3"

type Box struct {
	LowerBound, UpperBound vec3.T
}

func (this Box) BoundTriangle(tri Triangle) {
	this.LowerBound = vec3.T{
		min(tri[0][0], tri[0][1], tri[0][2]),
		min(tri[1][0], tri[1][1], tri[1][2]),
		min(tri[2][0], tri[2][1], tri[2][2]),
	}

	this.UpperBound = vec3.T{
		max(tri[0][0], tri[0][1], tri[0][2]),
		max(tri[1][0], tri[1][1], tri[1][2]),
		max(tri[2][0], tri[2][1], tri[2][2]),
	}
}

func max(a, b, c float32) float32 {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}

func min(a, b, c float32) float32 {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
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
