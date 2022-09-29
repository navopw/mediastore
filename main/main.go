package main

import (
	"encoding/json"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/labstack/echo/v5"
	"github.com/navopw/mediastore/env"
	_ "github.com/navopw/mediastore/migrations"
	"github.com/navopw/mediastore/routes"
	"github.com/navopw/mediastore/service"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	"net/http"
	"os"
)

func main() {
	// Env
	envErr := env.Initialize()
	if envErr != nil {
		panic(envErr)
	}

	// Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if !env.EnvironmentConfig.StructLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Libvips
	vips.Startup(nil)
	defer vips.Shutdown()

	// Exif parser
	exif.RegisterParsers(mknote.All...)

	// Pocketbase
	app := pocketbase.New()

	// S3
	fileSystem, err := filesystem.NewS3(
		env.EnvironmentConfig.MediaBucket,
		env.EnvironmentConfig.AwsRegion,
		env.EnvironmentConfig.S3Endpoint,
		env.EnvironmentConfig.AwsAccessKeyId,
		env.EnvironmentConfig.AwsSecretAccessKey,
		true,
	)
	if err != nil {
		panic(err)
	}

	service.FileSystem = fileSystem

	// Hook
	RegisterHooks(app)

	// Routes
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Routes
		middlewares := []echo.MiddlewareFunc{
			apis.RequireAdminOrUserAuth(),
		}
		e.Router.Add(http.MethodPost, "/api/media", routes.MediaCreateHandler, middlewares...)
		e.Router.Add(http.MethodGet, "/api/media", routes.MediaGetHandler, middlewares...)
		e.Router.Add(http.MethodDelete, "/api/media", routes.MediaDeleteHandler, middlewares...)
		e.Router.Add(http.MethodGet, "/api/media/totalsize", routes.MediaTotalSizeHandler, middlewares...)
		e.Router.Add(http.MethodGet, "/api/media/preview", routes.MediaPreviewHandler, middlewares...)
		e.Router.Add(http.MethodGet, "/api/media/list", routes.MediaListHandler, middlewares...)
		return nil
	})

	// Start
	service.App = app
	if err := app.Start(); err != nil {
		panic(err)
	}
}

func RegisterHooks(pb *pocketbase.PocketBase) {
	pb.OnModelBeforeCreate().Add(func(e *core.ModelEvent) error {
		marshal, err := json.Marshal(e.Model)
		if err != nil {
			return err
		}

		log.Print("Created model ", e.Model.TableName(), " "+string(marshal))

		return nil
	})

	pb.OnModelBeforeDelete().Add(func(e *core.ModelEvent) error {
		marshal, err := json.Marshal(e.Model)
		if err != nil {
			return err
		}

		log.Print("Deleted model ", e.Model.TableName(), " "+string(marshal))

		return nil
	})
}
