package env

import (
	carEnv "github.com/caarlos0/env"
)

type Environment struct {
	AwsAccessKeyId     string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion          string `env:"AWS_REGION"`
	S3Endpoint         string `env:"S3_ENDPOINT"`

	MediaBucket string `env:"MEDIA_BUCKET"`

	StructLog bool `env:"STRUCT_LOG"`
}

var EnvironmentConfig Environment

func Initialize() error {
	EnvironmentConfig = Environment{}
	err := carEnv.Parse(&EnvironmentConfig)
	if err != nil {
		return err
	}
	return nil
}
