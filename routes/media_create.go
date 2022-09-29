package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/navopw/photostore/model"
	"github.com/navopw/photostore/service"
	"github.com/navopw/photostore/util"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func MediaCreateHandler(c echo.Context) error {
	user, _ := c.Get(apis.ContextUserKey).(*models.User)

	// Form file
	formFile, formFileError := c.FormFile("file")

	if formFileError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "missing file",
			Error:   formFileError.Error(),
		})
	}

	// Open file
	file, formFileOpenError := formFile.Open()
	if formFileOpenError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to open file",
			Error:   formFileOpenError.Error(),
		})
	}

	mime := formFile.Header.Get("Content-Type")

	// File get bytes
	byteArray, byteReadError := io.ReadAll(file)
	if byteReadError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to read file",
			Error:   byteReadError.Error(),
		})
	}

	// SHA256 hash
	hash := sha256.Sum256(byteArray)
	hashString := hex.EncodeToString(hash[:])

	// Duplicate file check
	var mediaByHash []model.Media
	err := service.App.Dao().ModelQuery(&model.Media{}).Where(dbx.Like("hash", hashString)).Limit(1).All(&mediaByHash)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to check for duplicate file",
			Error:   err.Error(),
		})
	}

	if len(mediaByHash) > 0 {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "duplicate file",
		})
	}

	// Compress file
	log.Print("Compressing file...")
	compressedBytes := util.Compress(byteArray, mime)

	// Construct Media record
	media := model.Media{
		UserId:   user.Id,
		Hash:     hashString,
		Name:     formFile.Filename,
		Mime:     mime,
		FileSize: 1337,
	}
	media.Id = uuid.NewString()
	media.MarkAsNew()

	// Upload to S3
	log.Print("Uploading...")
	uploadError := service.FileSystem.Upload(compressedBytes, media.Id)
	if uploadError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to upload file",
			Error:   uploadError.Error(),
		})
	}

	// Preview
	log.Print("Extracting preview...")
	previewBytes, previewExtractError := util.ExtractPreviewByMime(compressedBytes, mime)
	if previewExtractError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to extract preview",
			Error:   previewExtractError.Error(),
		})
	}

	// Upload video preview
	log.Print("Uploading preview...")
	previewId := uuid.NewString()
	previewUploadError := service.FileSystem.Upload(previewBytes, previewId)
	if previewUploadError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to upload preview",
			Error:   previewUploadError.Error(),
		})
	}

	// Assign preview to media
	media.PreviewID = previewId

	// Save exif data into media
	log.Print("Extracting exif...")
	exifError := util.ExtractExifData(byteArray, &media)
	if exifError == nil {
		log.Print("Exif data found")
	}

	// Create Media
	log.Print("Creating media...")

	saveError := service.App.Dao().Save(&media)
	if saveError != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to save media",
			Error:   saveError.Error(),
		})
	}

	// Response
	return c.JSON(http.StatusOK, media)
}
