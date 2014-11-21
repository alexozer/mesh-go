package meshutil

type ArrayBuffer []Triangle

func (this *ArrayBuffer) NumTriangles() int {
	return len(*this)
}

func (this *ArrayBuffer) read() <-chan Triangle {
	triChan := make(chan Triangle)
	go func() {
		for _, tri := range *this {
			triChan <- tri
		}

		close(triChan)
	}()

	return triChan
}

func (this *ArrayBuffer) ConvertFrom(mesh Mesh) {
	*this = make(ArrayBuffer, 0, mesh.NumTriangles()*3)

	for tri := range mesh.read() {
		*this = append(*this, tri)
	}
}
