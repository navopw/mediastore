package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/navopw/mediastore/model"
	"github.com/navopw/mediastore/service"
	"github.com/navopw/mediastore/util"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"net/http"
)

type FilesystemSizeResponse struct {
	Size float64 `json:"size"`
}

func MediaTotalSizeHandler(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Find all media
	var media []model.Media
	err := service.App.Dao().
		ModelQuery(&model.Media{}).
		Select("file_size").
		Where(dbx.HashExp{"user_id": user.Id}).
		All(&media)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to get media from database",
			Error:   err.Error(),
		})
	}

	// Get total size from actual files
	var size float64

	for _, m := range media {
		size += m.FileSize
	}

	return c.JSON(http.StatusOK, FilesystemSizeResponse{
		Size: size,
	})
}
