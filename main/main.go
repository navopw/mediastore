package main

import (
	"encoding/json"
	"github.com/labstack/echo/v5"
	"github.com/navopw/mediastore/env"
	_ "github.com/navopw/mediastore/migrations"
	"github.com/navopw/mediastore/routes"
	"github.com/navopw/mediastore/service"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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

	// Exif parser
	exif.RegisterParsers(mknote.All...)

	// Pocketbase
	app := pocketbase.New()

	// S3
	ConfigureS3ByEnv(app)

	fileSystem, fileSystemError := app.NewFilesystem()
	if fileSystemError != nil {
		panic(fileSystemError)
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

func ConfigureS3ByEnv(pb *pocketbase.PocketBase) {
	s3Config := pb.Settings().S3
	s3Config.Enabled = true
	s3Config.Bucket = env.EnvironmentConfig.MediaBucket
	s3Config.Region = env.EnvironmentConfig.AwsRegion
	s3Config.Endpoint = env.EnvironmentConfig.S3Endpoint
	s3Config.AccessKey = env.EnvironmentConfig.AwsAccessKeyId
	s3Config.Secret = env.EnvironmentConfig.AwsSecretAccessKey
	s3Config.ForcePathStyle = true
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
}
