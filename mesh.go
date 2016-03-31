package mesh

import "github.com/ungerik/go3d/float64/vec3"

type Triangle [3]vec3.T

type Mesh interface {
	NumTriangles() int
	read() <-chan Triangle
	ConvertFrom(Mesh)
}
