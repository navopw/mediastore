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

func MediaGetHandler(c echo.Context) error {
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

	// Serve
	err = service.FileSystem.Serve(c.Response(), id, media.Name)

	// Set custom content type
	c.Response().Header().Set("Content-Type", media.Mime)

	// Failed to serve
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Message: "failed to serve media",
			Error:   err.Error(),
		})
	}

	return nil
}
