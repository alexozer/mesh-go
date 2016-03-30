package mesh

import (
	"encoding/binary"
	"errors"
	"os"

	"github.com/ungerik/go3d/float64/vec3"
)

const (
	uint16Size  = 2
	uint32Size  = 4
	float32Size = 4
)

var (
	errFormat = errors.New("wrong file format")
	errSize   = errors.New("reported face count does not match actual")
)

func IsFormat(err error) bool {
	return err == errFormat
}

func IsSize(err error) bool {
	return err == errSize
}

const (
	asciiId   = "solid"
	headerLen = 80
)

// StlFile reads from an existing .stl file or creates a new one if it doesn't exist.
type StlFile struct {
	path         string
	numTriangles int
}

func NewStlFile(filepath string) (*StlFile, error) {
	stlFile := StlFile{path: filepath}

	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return &stlFile, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stlFileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stlFileStat.Size() == 0 {
		return &stlFile, nil
	}

	if binary, err := stlFile.isBinary(file); err != nil {
		return nil, err
	} else if !binary {
		return nil, errFormat
	}

	if err = stlFile.readTriangleCount(file); err != nil {
		return nil, err
	}

	return &stlFile, nil
}

func (this *StlFile) isBinary(stlFile *os.File) (result bool, err error) {
	if _, err = stlFile.Seek(0, 0); err != nil {
		return
	}

	asciiText := make([]byte, len(asciiId))
	if _, err = stlFile.Read(asciiText); err != nil {
		return
	}

	return string(asciiText) != asciiId, nil
}

func (this *StlFile) readTriangleCount(stlFile *os.File) (err error) {
	_, err = stlFile.Seek(headerLen, 0)
	if err != nil {
		return
	}

	var numTriangles uint32
	err = binary.Read(stlFile, binary.LittleEndian, &numTriangles)
	if err != nil {
		return
	}
	this.numTriangles = int(numTriangles)

	var meshSize int64 = int64(this.numTriangles) * (3*4*float32Size + uint16Size)
	var fileSize int64 = headerLen + uint32Size + meshSize

	stats, err := stlFile.Stat()
	if err != nil {
		return
	}
	if stats.Size() != fileSize {
		return errSize
	}

	return
}

func (this *StlFile) read() <-chan Triangle {
	triChan := make(chan Triangle)
	go func() {
		file, err := os.Open(this.path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.Seek(headerLen+uint32Size, 0)

		var incomingTri stlTriangle

		var currTriangle int
		for ; currTriangle < this.numTriangles; currTriangle++ {
			binary.Read(file, binary.LittleEndian, &incomingTri)
			triChan <- incomingTri.Tri
		}

		close(triChan)
	}()

	return triChan
}

type stlTriangle struct {
	_   vec3.T // Normal
	Tri Triangle
	_   uint16 // Extra info
}

func (this *StlFile) ConvertFrom(mesh Mesh) {
	this.numTriangles = mesh.NumTriangles()

	err := os.Remove(this.path)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.Create(this.path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	header := make([]byte, headerLen)
	for i := range header {
		header[i] = 0x20
	}
	_, err = file.Write(header)
	if err != nil {
		panic(err)
	}

	err = binary.Write(file, binary.LittleEndian, uint32(this.numTriangles))
	if err != nil {
		panic(err)
	}

	var stlTri stlTriangle
	for tri := range mesh.read() {
		stlTri.Tri = tri
		binary.Write(file, binary.LittleEndian, stlTri)
	}
}

func (this *StlFile) NumTriangles() int {
	return this.numTriangles
}
