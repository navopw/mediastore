package mediahandler

import (
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"os"
	"path"
	"strings"
)

func GetExtensionFromMimeType(mime string) string {
	return mime[strings.LastIndex(mime, "/")+1:]
}

func Compress(byteArray []byte, mime string) []byte {
	if IsImage(mime) {
		return CompressImage(byteArray, mime)
	}

	if IsVideo(mime) {
		return CompressVideo(byteArray, mime)
	}

	// Dont compress
	log.Print("No compression for mime type", mime)
	return byteArray
}

func CompressImage(byteArray []byte, mime string) []byte {
	image, err := vips.NewImageFromBuffer(byteArray)
	if err != nil {
		log.Print("Failed to load image", err.Error())
		return byteArray
	}

	compressedBytes, _, err := image.ExportNative()
	if err != nil {
		return compressedBytes
	}

	_ = image.RemoveMetadata()

	log.Print("Compressed image from ", len(byteArray), " to ", len(compressedBytes))
	return compressedBytes
}

func CompressVideo(byteArray []byte, mime string) []byte {
	args := ffmpeg_go.KwArgs{
		"c:v":       "libx264",
		"profile:v": "main",
		"c:a":       "aac",
		"movflags":  "+faststart",
	}

	// Get temp folder
	tempFolder := os.TempDir()

	// Input file
	inputFilePath := path.Join(tempFolder, uuid.NewString()+"."+GetExtensionFromMimeType(mime))

	// Write input file
	err := os.WriteFile(inputFilePath, byteArray, 0644)
	if err != nil {
		log.Print("Failed to write video to temp folder", err.Error())
		return byteArray
	}

	// Output file path
	outputFile := path.Join(tempFolder, uuid.NewString()+"."+GetExtensionFromMimeType(mime))

	// Compress
	ffmpegError := ffmpeg_go.Input(inputFilePath).Output(outputFile, args).Run()
	if ffmpegError != nil {
		log.Print("Failed to compress video", ffmpegError.Error())
		return byteArray
	}

	// Read output file
	outputFileBytes, outputFileReadError := os.ReadFile(outputFile)
	if outputFileReadError != nil {
		log.Print("Failed to read compressed video", outputFileReadError.Error())
		return byteArray
	}

	// TODO: Error Handling
	// Delete input file & output file
	_ = os.Remove(inputFilePath)
	_ = os.Remove(outputFile)

	// Return compressed file
	log.Print("Compressed video from ", len(byteArray), " to ", len(outputFileBytes))
	return outputFileBytes
}
