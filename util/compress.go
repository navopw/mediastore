package util

import (
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/rs/zerolog/log"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"os"
	"path"
	"strings"
)

var imageMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
	"image/heic",
}

var videoMimeTypes = []string{
	"video/mp4",
	"video/webm",
	"video/ogg",
	"video/quicktime",
	"video/mov",
}

var pdfMimeType = "application/pdf"

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

	if mime == pdfMimeType {
		return CompressPdf(byteArray, mime)
	}

	// Dont compress
	log.Print("No compression for mime type", mime)
	return byteArray
}

func CompressImage(byteArray []byte, mime string) []byte {
	process, err := bimg.NewImage(byteArray).Process(bimg.Options{})
	if err != nil {
		return byteArray
	}

	log.Print("Compressed image from ", len(byteArray), " to ", len(process))
	return process
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

func CompressPdf(byteArray []byte, mime string) []byte {
	log.Print("PDF compression not implemented")
	return byteArray
}

func IsImage(mime string) bool {
	for _, imageMimeType := range imageMimeTypes {
		if imageMimeType == mime {
			return true
		}
	}
	return false
}

func IsVideo(mime string) bool {
	for _, videoMimeType := range videoMimeTypes {
		if videoMimeType == mime {
			return true
		}
	}
	return false
}
