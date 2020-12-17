package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const hankyungMarketInsightCommonURL = "http://marketinsight.hankyung.com/apps.free/free.news.list?category=IB_FREE"
const hankyungMarketInsightItemURL = "http://marketinsight.hankyung.com"

type HankyungMarketInsight struct{}

func (c HankyungMarketInsight) GetName() string {
	return "HankyungMarketInsight"
}

func (c HankyungMarketInsight) GetGroup() string {
	return "11"
}

func (c HankyungMarketInsight) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(hankyungMarketInsightCommonURL)
	if err != nil {
		return result, err
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// HTML Load
	html, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		return result, err
	}

	// Parsing
	wrapper := html.Find("ul#newsList")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a.txt")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := hankyungMarketInsightItemURL + href

		date := sel.Find("p.date").Text()
		title, _ := utils.ReadCP949(aTag.Find("strong > em").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c HankyungMarketInsight) GetContents(item *models.NewsItem) error {
	// Request
	req, err := http.Get(item.URL)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// HTML Load
	html, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		return err
	}

	// Parsing
	wrapper := html.Find("div#newsView")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			_contents, err := utils.ReadCP949(sel.Text())
			if err != nil {
				return
			}
			contents += (utils.TrimAll(_contents) + " ")
		}
	})

	item.Contents = contents
	return nil
}
