package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"path"
)

func ExtractPreviewFromVideo(byteArray []byte, mimeType string) ([]byte, error) {
	// Get temp folder
	tempFolder := os.TempDir()

	// Input path
	inputPath := path.Join(tempFolder, uuid.NewString()+"."+GetExtensionFromMimeType(mimeType))

	// Save input file
	saveError := os.WriteFile(inputPath, byteArray, 0644)
	if saveError != nil {
		return nil, saveError
	}

	frameBuffer, extractFrameError := ExtractFrameFromVideo(inputPath, 2)
	if extractFrameError != nil {
		return nil, extractFrameError
	}

	// Delete input file
	removeError := os.Remove(inputPath)
	if removeError != nil {
		return nil, removeError
	}

	return io.ReadAll(frameBuffer)
}

func ExtractFrameFromVideo(inFileName string, frameNum int) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(inFileName).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ExtractPreviewByMime(byteArray []byte, mime string) ([]byte, error) {
	// Video preview
	// Get first frame of video as jpg
	if IsVideo(mime) {
		videoPreviewBytes, videoExtractError := ExtractPreviewFromVideo(byteArray, mime)
		if videoExtractError != nil {
			return nil, videoExtractError
		}

		// Crop
		previewImage := bimg.NewImage(videoPreviewBytes)
		return previewImage.Process(bimg.Options{
			Width:       400,
			Height:      400,
			Crop:        true,
			Quality:     40,
			Compression: 50,
		})
	}

	// Image preview
	// Low quality
	if IsImage(mime) {
		previewImage := bimg.NewImage(byteArray)

		return previewImage.Process(bimg.Options{
			Width:       400,
			Height:      400,
			Crop:        true,
			Quality:     40,
			Compression: 50,
		})
	}

	return nil, errors.New("No preview for mime type " + mime)
}
