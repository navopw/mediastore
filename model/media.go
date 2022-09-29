package model

import (
	"github.com/pocketbase/pocketbase/models"
)

type Media struct {
	models.BaseModel

	// User
	UserId string `db:"user_id" json:"user_id"`

	// Media data
	Hash     string  `db:"hash" json:"hash"`
	Name     string  `db:"name" json:"name"`
	Mime     string  `db:"mime" json:"mime"`
	FileSize float64 `db:"file_size" json:"file_size"`

	// Video Preview image (S3 key)
	PreviewID string `db:"preview_id" json:"preview_id"`

	// Exif
	CameraModel string  `db:"camera_model" json:"camera_model"`
	FocalLength string  `db:"focal_length" json:"focal_length"`
	ExifDate    string  `db:"exif_date" json:"exif_date"`
	Latitude    float64 `db:"latitude" json:"latitude"`
	Longitude   float64 `db:"longitude" json:"longitude"`
}

func (e *Media) TableName() string {
	return "media"
}
