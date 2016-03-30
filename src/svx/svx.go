package svx

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"

	. "github.com/alexozer/go-mesh"
	"github.com/ungerik/go3d/vec2"
)

var ErrEmptyMesh = errors.New("Cannot export empty mesh")

const dirMode = 0755

func Export(mesh ArrayBuffer, path string, author string, voxelSize float32) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Put all the files in a temporary directory
	buildDir := path + ".d"
	err = os.RemoveAll(buildDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.Mkdir(buildDir, dirMode)
	if err != nil {
		return err
	}
	//defer os.RemoveAll(buildDir)

	slicePath := buildDir + string(os.PathSeparator) + sliceDir
	err = os.Mkdir(slicePath, dirMode)
	if err != nil {
		return err
	}

	if mesh.NumTriangles() == 0 {
		return ErrEmptyMesh
	}

	boxedTris := boxTriangles(mesh)
	totalBox := totalBox(boxedTris)

	// The model is in mm, but voxelSize is in m, so convert to mm
	voxelSizeMM := voxelSize / mmSize

	var sliceCounter int
	for z := totalBox.LowerBound[2]; z <= totalBox.UpperBound[2]; z += voxelSizeMM {
		layer := newLayer(totalBox, voxelSizeMM)
		for _, tri := range boxedTris {
			if z > tri.LowerBound[2] && z < tri.UpperBound[2] {
				line := zPlane(z).intersectTriangle(tri.Triangle)
				if line == nil {
					continue
				}

				layer.addLine(*line)
			}
		}
		layer.fill()

		sliceFilePath := slicePath + string(os.PathSeparator)
		sliceFilePath += fmt.Sprintf(sliceFormat, sliceCounter)
		sliceFile, err := os.Create(sliceFilePath)
		if err != nil {
			return err
		}

		// DEBUG
		if sliceCounter == 480 {
			for _, line := range layer.planeLines {
				fmt.Printf("%v\n", line)
			}
		}

		png.Encode(sliceFile, layer.Img)
		sliceFile.Close()

		sliceCounter++
	}

	return nil
}

func boxTriangles(mesh ArrayBuffer) []BoxedTriangle {
	boxedTris := make([]BoxedTriangle, mesh.NumTriangles())

	for i := range mesh {
		tri := mesh[i]
		box := BoxTriangle(tri)
		boxedTris[i] = BoxedTriangle{&tri, box}
	}

	return boxedTris
}

func totalBox(boxedTris []BoxedTriangle) *Box {
	if len(boxedTris) == 0 {
		return nil
	}

	initBbox := *boxedTris[0].Box

	for _, btri := range boxedTris[1:] {
		initBbox.AddBox(btri.Box)
	}

	return &initBbox
}

type zPlane float32 // Z-intercept

type planeLine [2]vec2.T

// intersectTriangle returns the directed line segment of the intersection,
// or nil if no intersection is found.
// Triangles that don't cross the plane do not intersect
func (this zPlane) intersectTriangle(tri *Triangle) *planeLine {
	result := planeLine{}

	var interceptFound bool
	for i := 0; i < 3; i++ {
		nextI := (i + 1) % 3
		line := &Line{tri[i], tri[nextI]}
		if !this.intersectsLine(line) {
			continue
		}
		interceptFound = true
		pt := this.intersectLine(line)

		if tri[nextI][2] > tri[i][2] {
			// Line "emerges" from plane
			result[1] = *pt
		} else {
			// Line "submerges" below the plane
			result[0] = *pt
		}
	}

	if !interceptFound {
		return nil
	}
	return &result
}

func (this zPlane) intersectLine(line *Line) *vec2.T {
	z := float32(this)
	scale := (z - line[0][2]) / (line[1][2] - line[0][2])
	x := scale*(line[1][0]-line[0][0]) + line[0][0]
	y := scale*(line[1][1]-line[0][1]) + line[0][1]

	return &vec2.T{x, y}
}

func (this zPlane) intersectsLine(line *Line) bool {
	zVal := float32(this)
	return (zVal > line[0][2] && zVal < line[1][2]) ||
		(zVal > line[1][2] && zVal < line[0][2])
}

func (this *planeLine) intersectsHorizLine(y, voxelSizeMM float32) bool {
	return (y >= this[0][1] && y <= this[1][1]) ||
		(y >= this[1][1] && y <= this[0][1])
}

func (this *planeLine) intersectHorizLine(y float32) (x float32) {
	deltaX := this[1][0] - this[0][0]
	deltaY := this[1][1] - this[0][1]
	dy := this[0][1] - y

	return this[0][0] - dy*deltaX/deltaY
}

func (this *planeLine) pointsUp() bool {
	return this[1][1] > this[0][1]
}

type layer struct {
	Img        *image.NRGBA
	VoxelSize  float32
	planeLines []planeLine
}

const mmSize = 1e-3

func newLayer(bounds *Box, voxelSize float32) *layer {
	return &layer{
		Img: image.NewNRGBA(image.Rect(
			int(bounds.LowerBound[0]/voxelSize),
			int(bounds.LowerBound[1]/voxelSize),
			int(bounds.UpperBound[0]/voxelSize),
			int(bounds.UpperBound[1]/voxelSize),
		)),
		VoxelSize:  voxelSize,
		planeLines: make([]planeLine, 0),
	}
}

func (this *layer) addLine(line planeLine) {
	this.planeLines = append(this.planeLines, line)
}

type intercept struct {
	X        float32
	PointsUp bool
}

type intercepts []intercept

func (this intercepts) Len() int {
	return len(this)
}

func (this intercepts) Less(i, j int) bool {
	return this[i].X < this[j].X
}

func (this intercepts) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this *layer) fill() {
	bounds := this.Img.Bounds()
	for imgY := bounds.Min.Y; imgY <= bounds.Max.Y; imgY++ {
		planeY := float32(imgY) * this.VoxelSize
		intercepts := make(intercepts, 0)

		for _, line := range this.planeLines {
			if !line.intersectsHorizLine(planeY, this.VoxelSize) {
				continue
			}

			intercepts = append(intercepts, intercept{
				line.intersectHorizLine(planeY),
				line.pointsUp(),
			})
		}
		sort.Sort(intercepts)

		var depth int
		for i, intercept := range intercepts {
			if !intercept.PointsUp {
				depth++
			} else {
				depth--
			}
			if i+1 < len(intercepts) && depth > 0 {
				imgX0, imgX1 := int(intercept.X/this.VoxelSize), int(intercepts[i+1].X/this.VoxelSize)
				this.fillStrip(imgX0, imgX1, imgY)
			}
		}
	}
}

func (this *layer) fillStrip(imgX0, imgX1, imgY int) {
	realImgY := this.Img.Rect.Max.Y - imgY + this.Img.Rect.Min.Y
	for x := imgX0; x <= imgX1; x++ {
		this.Img.Set(x, realImgY, color.NRGBA{255, 255, 255, 255})
	}
}
