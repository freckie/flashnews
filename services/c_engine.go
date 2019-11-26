package services

import (
	"flashnews/crawlers"
)

// CEngine : Crawling Engine
type CEngine struct {
	TG       *TGEngine
	Crawlers []*crawlers.Crawler
}
