package mesh

type Mesh interface {
	NumTriangles() int
	read() <-chan Triangle
	ConvertFrom(Mesh)
}
