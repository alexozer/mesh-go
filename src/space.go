package mesh

import (
	"math"

	"github.com/ungerik/go3d/float64/vec3"
)

type Box struct {
	LowerBound, UpperBound vec3.T
}

func (this *Box) Center() vec3.T {
	return vec3.T{
		(this.LowerBound[0] + this.UpperBound[0]) / 2,
		(this.LowerBound[1] + this.UpperBound[1]) / 2,
		(this.LowerBound[2] + this.UpperBound[2]) / 2,
	}
}

func (this *Box) ExpandToCube() *Box {
	var boxDimensions [3]float64
	for i := range boxDimensions {
		boxDimensions[i] = this.UpperBound[i] - this.LowerBound[i]
	}

	radius := math.Max(math.Max(boxDimensions[0], boxDimensions[1]), boxDimensions[2]) / 2
	center := this.Center()

	*this = Box{
		vec3.T{center[0] - radius, center[1] - radius, center[2] - radius},
		vec3.T{center[0] + radius, center[1] + radius, center[2] + radius},
	}

	return this
}

func BoxTriangle(tri Triangle) *Box {
	return &Box{
		vec3.T{
			math.Min(math.Min(tri[0][0], tri[1][0]), tri[2][0]),
			math.Min(math.Min(tri[0][1], tri[1][1]), tri[2][1]),
			math.Min(math.Min(tri[0][2], tri[1][2]), tri[2][2]),
		},
		vec3.T{
			math.Max(math.Max(tri[0][0], tri[1][0]), tri[2][0]),
			math.Max(math.Max(tri[0][1], tri[1][1]), tri[2][1]),
			math.Max(math.Max(tri[0][2], tri[1][2]), tri[2][2]),
		},
	}
}

func BoxTriangles(tris ...Triangle) *Box {
	if len(tris) == 0 {
		return nil
	}

	box := BoxTriangle(tris[0])
	for _, tri := range tris[1:] {
		for vert := 0; vert < 3; vert++ {
			for compon := 0; compon < 3; compon++ {
				val := tri[vert][compon]
				if val < box.LowerBound[compon] {
					box.LowerBound[compon] = val
				} else if val > box.UpperBound[compon] {
					box.UpperBound[compon] = val
				}
			}
		}
	}

	return box
}

func (this *Box) Intersects(other *Box) bool {
	if (this.LowerBound[0] > other.UpperBound[0]) ||
		(this.UpperBound[0] < other.LowerBound[0]) ||
		(this.LowerBound[1] > other.UpperBound[1]) ||
		(this.UpperBound[1] < other.LowerBound[1]) ||
		(this.LowerBound[2] > other.UpperBound[2]) ||
		(this.UpperBound[2] < other.LowerBound[2]) {
		return false
	}

	return true
}

func (this *Box) AddBox(other *Box) {
	this.LowerBound = vec3.T{
		math.Min(this.LowerBound[0], other.LowerBound[0]),
		math.Min(this.LowerBound[1], other.LowerBound[1]),
		math.Min(this.LowerBound[2], other.LowerBound[2]),
	}

	this.UpperBound = vec3.T{
		math.Max(this.UpperBound[0], other.UpperBound[0]),
		math.Max(this.UpperBound[1], other.UpperBound[1]),
		math.Max(this.UpperBound[2], other.UpperBound[2]),
	}
}

func (this *Box) IntersectsXY(z float64) bool {
	return z > this.LowerBound[2] && z < this.UpperBound[2]
}

type BoxedTriangle struct {
	*Triangle
	*Box
}

const octreeMaxTriangles = 10

type Octree struct {
	*Box
	level     uint
	isLeaf    bool
	children  [8]*Octree
	triangles []BoxedTriangle
}

func NewOctree(levels uint, box *Box) *Octree {
	if levels == 0 {
		return nil
	}

	return &Octree{box, levels - 1, true, [8]*Octree{}, make([]BoxedTriangle, 0)}
}

func (this *Octree) Insert(tri *BoxedTriangle) {
	if this.level == 0 {
		this.triangles = append(this.triangles, *tri)
		return
	}

	if this.isLeaf {
		this.triangles = append(this.triangles, *tri)
		if len(this.triangles) > octreeMaxTriangles {
			this.isLeaf = false
			for _, childTri := range this.triangles {
				this.Insert(&childTri)
			}
			this.triangles = nil
		}
		return
	}

	if this.children[0] == nil {
		this.initChildren()
	}

	for _, node := range this.children {
		if node.Intersects(tri.Box) {
			node.Insert(tri)
		}
	}

	return
}

func (this *Octree) initChildren() {
	center := this.Center()
	for i := range this.children {
		child := NewOctree(this.level, new(Box))
		this.children[i] = child

		/* For a three-bit number:
		if bit == 0, octant in "left" part of axis
		if bit == 1, then octant in "right" part of axis
		bit #2: x-axis
		bit #1: y-axis
		bit #0: z-axis
		*/
		for i := uint(0); i < 3; i++ {
			if i&(4>>i) == 0 {

				child.LowerBound[i] = this.LowerBound[i]
				child.UpperBound[i] = center[i]
			} else {
				child.LowerBound[i] = center[i]
				child.UpperBound[i] = this.UpperBound[i]
			}
		}
	}
}
