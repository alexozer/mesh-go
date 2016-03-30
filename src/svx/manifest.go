package svx

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type manifest struct {
	XMLName           xml.Name `xml:"grid"`
	Version           string   `xml:"version,attr"`
	GridSizeX         int      `xml:"gridSizeX,attr"`
	GridSizeY         int      `xml:"gridSizeY,attr"`
	GridSizeZ         int      `xml:"gridSizeZ,attr"`
	VoxelSize         float32  `xml:"voxelSize,attr"`
	SubvoxelBits      int      `xml:"subvoxelBits,attr"`
	SlicesOrientation string   `xml:"slicesOrientation,attr"`

	Channel  channel         `xml:"channels>channel"`
	Material material        `xml:"materials>material"`
	Metadata []metadataEntry `xml:"metadata>entry"`
}

var defaultManifest = manifest{
	Version:           "1.0",
	SubvoxelBits:      8,
	SlicesOrientation: "Z",
}

const manifestHeader = `<?xml version="1.0"?>` + "\n\n"

func newManifest(author string, gridSizeX, gridSizeY, gridSizeZ int, voxelSize float32) *manifest {
	m := defaultManifest

	m.GridSizeX = gridSizeX
	m.GridSizeY = gridSizeY
	m.GridSizeZ = gridSizeZ
	m.VoxelSize = voxelSize

	m.Channel = defaultChannel
	m.Material = defaultMaterial
	m.Metadata = newMetadata(author)

	return &m
}

type channel struct {
	XMLName xml.Name `xml:"channel"`
	Type    string   `xml:"type,attr"`
	Bits    int      `xml:"bits,attr"`
	Slices  string   `xml:"slices,attr"`
}

const sliceDir = "density"
const sliceFormat = `slice%d.png`

var defaultChannel = channel{
	Type:   "DENSITY",
	Bits:   8,
	Slices: sliceDir + string(os.PathSeparator) + sliceFormat,
}

type material struct {
	XMLName xml.Name `xml:"material"`
	Id      int      `xml:"id,attr"`
	Urn     string   `xml:"urn,attr"`
}

var defaultMaterial = material{
	Id:  1,
	Urn: "urn:shapeways:materials/1",
}

type metadataEntry struct {
	XMLName xml.Name `xml:"entry"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

func newMetadata(author string) []metadataEntry {
	year, month, day := time.Now().Date()
	dateStr := fmt.Sprintf("%d/%d/%d", year, month, day)

	return []metadataEntry{
		{
			Key:   "author",
			Value: author,
		}, {
			Key:   "creationDate",
			Value: dateStr,
		},
	}
}

func (this *manifest) Export(filepath string) error {
	err := os.Remove(filepath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	file.WriteString(manifestHeader)

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	err = encoder.Encode(this)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
