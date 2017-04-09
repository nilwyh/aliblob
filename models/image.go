package models

import (
	"time"
)

type Image struct {
	CreatedTime time.Time
	Title       string
	Description string
	SoftTags    []string
	Src         string
	Domain      string
	SHA256      string
	Group       string
	IsRaw bool
	RawId string
	Thumb string
	GoodToShow bool
	Format string
}