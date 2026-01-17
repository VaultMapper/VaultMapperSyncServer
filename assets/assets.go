package assets

import (
	_ "embed"
	"strconv"
	"strings"
)

type Meta struct {
	Width     int
	Height    int
	FrameSize int
	Frames    int
	FPS       int
}

type Video struct {
	Data []byte
	Meta
}

type FrameData []byte

func (v *Video) GetFrameData(frameIndex int) (frameData FrameData) {
	bitStart := v.FrameSize * frameIndex
	byteIndex := bitStart / 8
	bitOffset := bitStart % 8

	// Calculate how many bytes needed
	numBytes := (v.FrameSize + bitOffset + 7) / 8
	result := make([]byte, (v.FrameSize+7)/8)

	// Copy and shift the bytes
	for i := 0; i < numBytes && byteIndex+i < len(v.Data); i++ {
		result[i] |= v.Data[byteIndex+i] << bitOffset
		if bitOffset > 0 && byteIndex+i+1 < len(v.Data) {
			result[i] |= v.Data[byteIndex+i+1] >> (8 - bitOffset)
		}
	}

	// Mask the last byte
	if remainder := v.FrameSize % 8; remainder != 0 {
		result[len(result)-1] &= byte(0xFF << (8 - remainder))
	}

	return result
}

func (f *FrameData) Iterate(callback func(bits byte)) {
	for i := 0; i < len(*f); i++ {
		callback((*f)[i])
	}
}

//go:embed bad_apple.bin
var badAppleData []byte

//go:embed bad_apple.bin.meta
var badAppleMetaString string
var badAppleMetaValues = strings.Split(badAppleMetaString, ",")
var badAppleMetaWidth, _ = strconv.Atoi(badAppleMetaValues[0])
var badAppleMetaHeight, _ = strconv.Atoi(badAppleMetaValues[1])
var badAppleMetaFrames, _ = strconv.Atoi(badAppleMetaValues[2])
var badAppleMetaFPS, _ = strconv.Atoi(badAppleMetaValues[3])

var badAppleMeta = Meta{
	badAppleMetaWidth,
	badAppleMetaHeight,
	badAppleMetaWidth * badAppleMetaHeight,
	badAppleMetaFrames,
	badAppleMetaFPS,
}

var BadApple = Video{
	badAppleData,
	badAppleMeta,
}
