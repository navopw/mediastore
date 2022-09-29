package routes

import (
	"errors"
	"github.com/labstack/echo/v5"
	"github.com/navopw/photostore/model"
	"github.com/navopw/photostore/service"
	"github.com/navopw/photostore/util"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"net/http"
)

func MediaPreviewHandler(c echo.Context) error {
	user, _ := c.Get(apis.ContextUserKey).(*models.User)
	id := c.QueryParam("id")

	// Validation
	if id == "" {
		return errors.New("missing query param id")
	}

	// Find media
	var media model.Media
	err := service.App.Dao().ModelQuery(&model.Media{}).Where(dbx.HashExp{"id": id}).One(&media)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.ErrorResponse{
			Message: "failed to find record",
			Error:   err.Error(),
		})
	}

	// User check
	if media.UserId != user.Id {
		return c.JSON(http.StatusUnauthorized, util.ErrorResponse{
			Message: "unauthorized",
		})
	}

	// Serve
	err = service.FileSystem.Serve(c.Response(), media.PreviewID, "record")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to serve preview",
			Error:   err.Error(),
		})
	}
	return nil
}
