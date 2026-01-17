package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Config struct {
	InputFile  string
	OutputFile string
	Width      int
	Height     int
	Threshold  uint
}

func main() {
	config := Config{}

	flag.StringVar(&config.InputFile, "input", "", "Input MP4 file path")
	flag.IntVar(&config.Width, "width", 32, "Output frame width")
	flag.IntVar(&config.Height, "height", 18, "Output frame height")
	flag.UintVar(&config.Threshold, "threshold", 128, "Black/white threshold (0-255)")
	flag.Parse()

	if config.InputFile == "" {
		log.Fatal("Input file is required")
	}

	if !strings.HasSuffix(config.InputFile, ".mp4") {
		log.Fatal("Input file must be .mp4")
	}

	config.OutputFile = strings.TrimSuffix(config.InputFile, ".mp4") + ".bin"

	if err := ProcessVideo(config); err != nil {
		log.Fatal(err)
	}
}

func ProcessVideo(config Config) error {
	// Extract frames and get dimensions using ffmpeg scaling
	// Let ffmpeg handle the initial extraction at target resolution
	stream := ffmpeg.Input(config.InputFile).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", config.Width, config.Height)})

	buf := &bytes.Buffer{}
	err := stream.Output("pipe:",
		ffmpeg.KwArgs{
			"format":  "rawvideo",
			"pix_fmt": "rgb24",
		}).
		WithOutput(buf).
		Run()

	if err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}

	frameSize := config.Width * config.Height * 3
	data := buf.Bytes()
	frameCount := len(data) / frameSize

	fmt.Printf("Output: %dx%d\n", config.Width, config.Height)
	fmt.Printf("Processing %d frames\n", frameCount)

	// Create output file
	outFile, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Process each frame
	for i := 0; i < frameCount; i++ {
		offset := i * frameSize
		frameData := data[offset : offset+frameSize]

		// Convert to binary and pack
		binaryData := convertToBinaryPacked(frameData, config.Width, config.Height, config.Threshold)

		if _, err := outFile.Write(binaryData); err != nil {
			return fmt.Errorf("failed to write frame %d: %w", i, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("Processed %d/%d frames\n", i+1, frameCount)
		}
	}

	outMetaFile, err := os.Create(config.OutputFile + ".meta")
	if err != nil {
		return fmt.Errorf("failed to create output meta file: %w", err)
	}
	defer outMetaFile.Close()

	fps, err := GetVideoFPS(config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to get video fps: %w", err)
	}

	outMetaFile.WriteString(fmt.Sprintf("%d,%d,%d,%d", config.Width, config.Height, frameCount, int(fps)))

	bytesPerFrame := (config.Width*config.Height + 7) / 8
	totalSize := bytesPerFrame * frameCount

	fmt.Printf("\nDone!\n")
	fmt.Printf("Output: %s\n", config.OutputFile)
	fmt.Printf("Frames: %d\n", frameCount)
	fmt.Printf("Bytes per frame: %d\n", bytesPerFrame)
	fmt.Printf("Total size: %d bytes (%.2f MB)\n", totalSize, float64(totalSize)/(1024*1024))

	return nil
}

func GetVideoFPS(inputFile string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=r_frame_rate",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputFile,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe error: %w", err)
	}

	// Parse the fraction (e.g., "30000/1001" or "30/1")
	fpsStr := strings.TrimSpace(string(output))
	parts := strings.Split(fpsStr, "/")

	if len(parts) != 2 {
		return 0, fmt.Errorf("unexpected fps format: %s", fpsStr)
	}

	numerator, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numerator: %w", err)
	}

	denominator, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse denominator: %w", err)
	}

	return numerator / denominator, nil
}

func convertToBinaryPacked(rawData []byte, width, height int, threshold uint) []byte {
	totalPixels := width * height
	numBytes := (totalPixels + 7) / 8
	packedData := make([]byte, numBytes)

	bitIndex := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 3
			r := rawData[idx]
			g := rawData[idx+1]
			b := rawData[idx+2]

			// Convert to grayscale
			gray := uint(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

			// Threshold to binary: >= threshold = 1 (white), < threshold = 0 (black)
			if gray >= threshold {
				byteIndex := bitIndex / 8
				bitPosition := 7 - (bitIndex % 8)
				packedData[byteIndex] |= (1 << bitPosition)
			}

			bitIndex++
		}
	}

	return packedData
}
