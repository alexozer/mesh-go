package mesh

import (
	"fmt"
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

	if !cubeArray1.Equals(&cubeArray2) {
		t.Fatal("Stl conversion error")
	}
}

func TestArrays(t *testing.T) {
	stl, err := NewStlFile(cubePath)
	if err != nil {
		t.Fatal(err)
	}

	abuf1 := ArrayBuffer{}
	abuf1.ConvertFrom(stl)
	abuf2 := ArrayBuffer{}
	abuf2.ConvertFrom(&abuf1)
	if !abuf1.Equals(&abuf2) {
		t.Fatal("ArrayBuffer conversions differ")
	}

	ibuf := IndexBuffer{}
	ibuf.ConvertFrom(&abuf1)
	abuf2.ConvertFrom(&ibuf)
	if !abuf1.Equals(&abuf2) {
		fmt.Println(ibuf)
		t.Fatal("IndexBuffer conversions differ")
	}
}

func (this *ArrayBuffer) Equals(other *ArrayBuffer) bool {
	if len(*this) != len(*other) {
		return false
	}

	for i, val := range *this {
		if val != (*other)[i] {
			return false
		}
	}

	return true
}
