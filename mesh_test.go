package meshutil

import (
	"io/ioutil"
	"os"
	"testing"
)

const cubePath = "resources/cube.stl"

func TestStl(t *testing.T) {
	stl1, err := NewStlFile(cubePath)
	if err != nil {
		t.Fatal(err)
	}

	cubeArray1 := ArrayBuffer{}
	cubeArray1.ConvertFrom(stl1)

	tmpfile, err := ioutil.TempFile("", "stlfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	stl2, err := NewStlFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	stl2.ConvertFrom(stl1)

	cubeArray2 := ArrayBuffer{}
	cubeArray2.ConvertFrom(stl2)

	if !equal(cubeArray1, cubeArray2) {
		t.Fatal("Stl conversion error")
	}
}

func equal(mesh1, mesh2 ArrayBuffer) bool {
	if len(mesh1) != len(mesh2) {
		return false
	}

	for i, val := range mesh1 {
		if val != mesh2[i] {
			return false
		}
	}

	return true
}
