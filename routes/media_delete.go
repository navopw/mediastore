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
)

func MediaDeleteHandler(c echo.Context) error {
	user, _ := c.Get(apis.ContextUserKey).(*models.User)
	id := c.QueryParam("id")

	// Validation
	if id == "" {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "missing id",
		})
	}

	// Find media
	var media model.Media
	err := service.App.Dao().ModelQuery(&model.Media{}).Where(dbx.HashExp{"id": id}).One(&media)

	// User check
	if media.UserId != user.Id {
		return c.JSON(http.StatusUnauthorized, util.ErrorResponse{
			Message: "unauthorized",
		})
	}

	// Delete media & preview
	err = service.FileSystem.Delete(media.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to delete media from storage",
			Error:   err.Error(),
		})
	}

	err = service.FileSystem.Delete(media.PreviewID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to delete media preview from storage",
			Error:   err.Error(),
		})
	}

	// Delete media from database
	err = service.App.Dao().Delete(&media)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to delete media from database",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, media)
}
