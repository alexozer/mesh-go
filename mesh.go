package meshutil

type vertex [3]float32

// TODO: remove NumVertices()
type Mesh interface {
	NumVertices() uint
	NumTriangles() uint
	read() <-chan vertex
	ConvertFrom(Mesh)
}

const (
	uint16Size  = 2
	uint32Size  = 4
	float32Size = 4
)
