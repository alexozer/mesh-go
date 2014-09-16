package meshutil

type ArrayBuffer []float32

func (buf *ArrayBuffer) NumVertices() uint {
	return uint(len(*buf)) / 3
}

func (buf *ArrayBuffer) NumTriangles() uint {
	return buf.NumVertices() / 3
}

func (buf *ArrayBuffer) read() <-chan vertex {
	verts := make(chan vertex)
	go func() {
		var vert vertex
		for bufPos := 0; bufPos < len(*buf); bufPos += len(vert) {
			copy(vert[:], (*buf)[bufPos:])
			verts <- vert
		}

		close(verts)
	}()

	return verts
}

func (buf *ArrayBuffer) ConvertFrom(mesh Mesh) {
	newBuf := make(ArrayBuffer, 0, mesh.NumVertices()*3)
	for vertex := range mesh.read() {
		newBuf = append(newBuf, vertex[:]...)
	}
	*buf = newBuf
}
