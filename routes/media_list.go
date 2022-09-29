package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/navopw/mediastore/model"
	"github.com/navopw/mediastore/service"
	"github.com/navopw/mediastore/util"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"net/http"
	"strings"
)

func MediaListHandler(c echo.Context) error {
	user, _ := c.Get(apis.ContextUserKey).(*models.User)

	// Filter
	videoFlag := c.QueryParam("video")
	imageFlag := c.QueryParam("image")

	var mimeFilter string
	if videoFlag == "true" {
		mimeFilter = "video"
	} else if imageFlag == "true" {
		mimeFilter = "image"
	}

	// Find all media
	var media []model.Media

	query := service.App.Dao().
		ModelQuery(&model.Media{}).
		Where(dbx.HashExp{"user_id": user.Id})

	// If filter is set, filter media
	if mimeFilter != "" {
		query = query.AndWhere(dbx.Like("mime", mimeFilter))
	}

	err := query.All(&media)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.ErrorResponse{
			Message: "failed to get media list",
			Error:   err.Error(),
		})
	}

	// List empty
	if len(media) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	// Return media list
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
