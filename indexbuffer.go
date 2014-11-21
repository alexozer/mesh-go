package meshutil

import "github.com/ungerik/go3d/vec3"

type Face [3]uint16

type IndexBuffer struct {
	Vertices []vec3.T
	Faces    []Face
}

func (this *IndexBuffer) NumTriangles() int {
	return len(this.Faces)
}

func (this *IndexBuffer) read() <-chan Triangle {
	triChan := make(chan Triangle)
	go func() {
		for _, face := range this.Faces {
			triChan <- Triangle{
				this.Vertices[face[0]],
				this.Vertices[face[1]],
				this.Vertices[face[2]],
			}
		}

		close(triChan)
	}()

	return triChan
}

func (this *IndexBuffer) ConvertFrom(mesh Mesh) {
	this.Vertices = make([]vec3.T, 0)
	this.Faces = make([]Face, mesh.NumTriangles())

	uniqueVertices := make(map[vec3.T]uint16)
	var currIndex uint16

	var face Face
	for tri := range mesh.read() {
		for i, vert := range tri[:] {
			index, exists := uniqueVertices[vert]

			if !exists {
				this.Vertices = append(this.Vertices, vert)
				uniqueVertices[vert] = currIndex

				index = currIndex
				currIndex++
			}
			face[i] = index
		}

		this.Faces = append(this.Faces, face)
	}
}
