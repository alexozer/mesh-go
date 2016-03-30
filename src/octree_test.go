package mesh

/*func TestOctree(t *testing.T) {
	stl, err := NewStlFile("resources/cube.stl")
	if err != nil {
		t.Fatal(err)
	}

	abuf := ArrayBuffer{}
	abuf.ConvertFrom(stl)

	bbox := BoxTriangles(abuf...).ExpandToCube()

	tree := NewOctree(8, bbox)

	for _, tri := range abuf {
		tree.Insert(NewBoxedTriangle(tri))
	}

	fmt.Println(tree)
}*/
