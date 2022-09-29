package mediahandler

import (
	"bytes"
	"github.com/navopw/mediastore/model"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/rwcarlsen/goexif/exif"
	"strings"
)

// TODO: Heic support
func ExtractExifData(byteArray []byte, media *model.Media) error {
	// Decode
	decoded, decodeError := exif.Decode(bytes.NewReader(byteArray))
	if decodeError != nil {
		return decodeError
	}

	// CamModel
	camModel, camModelError := decoded.Get(exif.Model)
	if camModelError == nil {
		media.CameraModel = strings.Trim(camModel.String(), "\"")
	}

	// FocalLength
	focalLength, focalLengthError := decoded.Get(exif.FocalLength)
	if focalLengthError == nil {
		media.FocalLength = strings.Trim(focalLength.String(), "\"")
	}

	// Date
	dateTime, dateTimeError := decoded.DateTime()
	if dateTimeError == nil {
		parsedDateTime, parseDateTimeError := types.ParseDateTime(dateTime)
		if parseDateTimeError == nil {
			media.ExifDate = parsedDateTime.String()
		}
	}

	// Latitude and Longitude
	lat, long, locationError := decoded.LatLong()
	if locationError == nil {
		media.Latitude = lat
		media.Longitude = long
	}

	return nil
}
