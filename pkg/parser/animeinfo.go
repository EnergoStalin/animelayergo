package parser

import (
	"time"
)

type AnimeInfo struct {
	Title        string
	DownloadLink string
	Seed         int
	Leech        int
	Size         int
	Uploader     string
	Date         time.Time
	Hash         string
	MagnetUri    string
}
