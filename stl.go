package meshutil

import (
	"encoding/binary"
	"errors"
	"os"
)

var (
	errFormat             = errors.New("wrong file format")
	errSize               = errors.New("reported face count does not match actual")
	errIncompleteTriangle = errors.New("vertices received not divisible by 3")
)

func IsFormat(err error) bool {
	return err == errFormat
}

func IsSize(err error) bool {
	return err == errSize
}

func IsIncompleteTriangle(err error) bool {
	return err == errIncompleteTriangle
}

const (
	asciiId   = "solid"
	headerLen = 80
)

// StlFile reads from an existing .stl file or creates a new one if it doesn't exist.
type StlFile struct {
	path      string
	triangles uint32
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

func (file *StlFile) isBinary(stlFile *os.File) (result bool, err error) {
	if _, err = stlFile.Seek(0, 0); err != nil {
		return
	}

	asciiText := make([]byte, len(asciiId))
	if _, err = stlFile.Read(asciiText); err != nil {
		return
	}

	return string(asciiText) != asciiId, nil
}

func (file *StlFile) readTriangleCount(stlFile *os.File) (err error) {
	_, err = stlFile.Seek(headerLen, 0)
	if err != nil {
		return
	}

	err = binary.Read(stlFile, binary.LittleEndian, &file.triangles)
	if err != nil {
		return
	}

	var meshSize int64 = int64(file.triangles) * (3*4*float32Size + uint16Size)
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

func (stl *StlFile) read() <-chan vertex {
	verts := make(chan vertex)
	go func() {
		file, err := os.Open(stl.path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.Seek(headerLen+uint32Size, 0)

		var triangle stlTriangle

		var currTriangle uint32
		for ; currTriangle < stl.triangles; currTriangle++ {
			binary.Read(file, binary.LittleEndian, &triangle)
			verts <- triangle.Vert1
			verts <- triangle.Vert2
			verts <- triangle.Vert3
		}

		close(verts)
	}()

	return verts
}

type stlTriangle struct {
	Normal [3]float32
	Vert1  vertex
	Vert2  vertex
	Vert3  vertex
	_      uint16
}

func (stl *StlFile) ConvertFrom(mesh Mesh) {
	stl.triangles = uint32(mesh.NumTriangles())

	err := os.Remove(stl.path)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.Create(stl.path)
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

	err = binary.Write(file, binary.LittleEndian, stl.triangles)
	if err != nil {
		panic(err)
	}

	vertChan := mesh.read()
	triangle := stlTriangle{}
	for triangle.Vert1 = range vertChan {
		var ok bool // Don't know why Go won't let me define 'ok' with := below
		triangle.Vert2, ok = <-vertChan
		var ok2 bool
		triangle.Vert3, ok2 = <-vertChan
		if !(ok && ok2) {
			panic(errIncompleteTriangle)
		}

		binary.Write(file, binary.LittleEndian, triangle)
	}
}

func (file *StlFile) NumTriangles() uint {
	return uint(file.triangles)
}

func (file *StlFile) NumVertices() uint {
	return file.NumTriangles() * 3
}
