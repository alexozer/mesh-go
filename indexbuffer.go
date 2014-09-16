package meshutil

type IndexBuffer struct {
	Vertices []float32
	Indices  []uint16
}

func (buf *IndexBuffer) NumVertices() uint {
	return uint(len(buf.Indices))
}

func (buf *IndexBuffer) NumTriangles() uint {
	return buf.NumVertices() / 3
}

func (buf *IndexBuffer) read() <-chan vertex {
	vertChan := make(chan vertex)
	go func() {
		var vert vertex
		for _, vertIndex := range buf.Indices {
			copy(vert[:], buf.Vertices[int(vertIndex)*len(vert):])
			vertChan <- vert
		}

		close(vertChan)
	}()

	return vertChan
}

func (buf *IndexBuffer) ConvertFrom(mesh Mesh) {
	uniqueVertices := make(map[vertex]uint16)
	var currIndex uint16

	for vert := range mesh.read() {
		index, exists := uniqueVertices[vert]
		if !exists {
			buf.Vertices = append(buf.Vertices, (vert)[:]...)
			uniqueVertices[vert] = currIndex

			index = currIndex
			currIndex++
		}
		buf.Indices = append(buf.Indices, index)
	}
}
