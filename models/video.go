package models

import "time"

type Video struct {
	CreatedTime time.Time
	Title       string
	Description string
	SoftTags    []string
	Src         string
	Domain      string
	SHA256      string
	Format      string
}