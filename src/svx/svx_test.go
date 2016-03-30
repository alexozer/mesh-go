package svx

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/alexozer/go-mesh"
	"github.com/ungerik/go3d/vec2"
	"github.com/ungerik/go3d/vec3"
)

func newAbuf(t *testing.T) mesh.ArrayBuffer {
	stl, err := mesh.NewStlFile("../resources/helix.stl")
	if err != nil {
		t.Fatal(err)
	}

	abuf := mesh.ArrayBuffer{}
	abuf.ConvertFrom(stl)
	return abuf
}

const (
	author     = "Alex Ozer"
	gridSize   = 100
	voxelSize  = 1e-4
	exportFile = "/tmp/export.svx"
	tmpImage   = "/tmp/test.png"
)

func TestManifest(t *testing.T) {
	manifest := newManifest(author, gridSize, gridSize, gridSize, voxelSize)
	manifest.Export("/tmp/manifest.xml")
}

func TestFill(t *testing.T) {
	layer := newLayer(&mesh.Box{vec3.T{0, 0, 0}, vec3.T{1, 1, 1}}, voxelSize)

	tri0, tri1, tri2 := vec2.T{0.25, 0.25}, vec2.T{0.75, 0.25}, vec2.T{0.5, 0.75}
	sq0, sq1, sq2, sq3 := vec2.T{0.3, 0.3}, vec2.T{0.7, 0.3}, vec2.T{0.7, 0.7}, vec2.T{0.3, 0.7}

	// Triangle
	layer.addLine(planeLine{tri0, tri2})
	layer.addLine(planeLine{tri2, tri1})
	layer.addLine(planeLine{tri1, tri0})

	// Square
	//layer.addLine(planeLine{sq3, sq2})
	//layer.addLine(planeLine{sq2, sq1})
	//layer.addLine(planeLine{sq1, sq0})
	//layer.addLine(planeLine{sq0, sq3})

	// Cut-out square
	layer.addLine(planeLine{sq0, sq1})
	layer.addLine(planeLine{sq1, sq2})
	layer.addLine(planeLine{sq2, sq3})
	layer.addLine(planeLine{sq3, sq0})

	layer.fill()

	os.Remove(tmpImage)
	file, err := os.Create(tmpImage)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	png.Encode(file, layer.Img)
}

func TestPng(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 500, 500))

	const testFile = "/tmp/thetestimage.png"
	os.Remove(testFile)
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func TestExport(t *testing.T) {
	abuf := newAbuf(t)

	path := "/tmp/test.svx"
	err := Export(abuf, path, author, voxelSize)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntersectTriangle(t *testing.T) {
	tri := &mesh.Triangle{
		vec3.T{0, 0, 0},
		vec3.T{0, 0, 100},
		vec3.T{0, 100, 0},
	}

	line := zPlane(50).intersectTriangle(tri)
	fmt.Println(line)
}

func TestIntersectTriangles(t *testing.T) {
	abuf := newAbuf(t)

	for i := range abuf {
		tri := abuf[i]
		fmt.Println(tri, zPlane(50).intersectTriangle(&tri))
	}
}

func TestIntersectLine(t *testing.T) {
	line := &mesh.Line{
		vec3.T{100, 0, 0},
		vec3.T{100, 100, 100},
	}

	fmt.Println(zPlane(50).intersectLine(line))
}
