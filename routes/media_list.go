package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/navopw/photostore/model"
	"github.com/navopw/photostore/service"
	"github.com/navopw/photostore/util"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"net/http"
	"strings"
)

func MediaListHandler(c echo.Context) error {
	user, _ := c.Get(apis.ContextUserKey).(*models.User)

	videoFlag := c.QueryParam("video")
	imageFlag := c.QueryParam("image")

	// Find all media
	var media []model.Media

	err := service.App.Dao().ModelQuery(&model.Media{}).Where(dbx.HashExp{"user_id": user.Id}).All(&media)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.ErrorResponse{
			Message: "failed to get media list",
			Error:   err.Error(),
		})
	}

	// Filter by video
	if videoFlag == "true" {
		media = filterMedia(media, "video")
	}

	// Filter by image
	if imageFlag == "true" {
		media = filterMedia(media, "image")
	}

	// List empty
	if len(media) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	return c.JSON(http.StatusOK, media)
}

// TODO: Database filter
func filterMedia(media []model.Media, prefix string) []model.Media {
	var filteredMedia []model.Media

	for _, m := range media {
		if strings.Contains(m.Mime, prefix) {
			filteredMedia = append(filteredMedia, m)
		}
	}

	return filteredMedia
}
