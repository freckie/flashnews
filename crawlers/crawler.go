package crawlers

import (
	"flashnews/models"
	"regexp"
)

type Crawler interface {
	GetName() string                               // Get Crawler Name
	GetGroup() string                              // Get Group Number ex. "4"
	GetList(number int) ([]models.NewsItem, error) // Crawl List of NewsItem
	GetContents(item *models.NewsItem) error       // Crawl Contents of Particular Item
}

func GetName(c Crawler) string {
	return c.GetName()
}

func GetGroup(c Crawler) string {
	return c.GetGroup()
}

func GetList(c Crawler, number int) ([]models.NewsItem, error) {
	return c.GetList(number)
}

func GetContents(c Crawler, item *models.NewsItem) error {
	return c.GetContents(item)
}

var /* const */ reForNumbers = regexp.MustCompile("[0-9]+")
