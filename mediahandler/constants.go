package mediahandler

var SupportedImageMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
	HeicMime,
}

var SupportedVideoMimeTypes = []string{
	"video/mp4",
	"video/webm",
	"video/ogg",
	"video/quicktime",
	"video/mov",
}

const HeicMime = "image/heic"

func IsImage(mime string) bool {
	for _, imageMimeType := range SupportedImageMimeTypes {
		if imageMimeType == mime {
			return true
		}
	}
	return false
}

func IsVideo(mime string) bool {
	for _, videoMimeType := range SupportedVideoMimeTypes {
		if videoMimeType == mime {
			return true
		}
	}
	return false
}

func IsSupported(mime string) bool {
	return IsImage(mime) || IsVideo(mime)
}
