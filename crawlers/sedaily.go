package crawlers

import (
	"fmt"
	"net/http"
	"regexp"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const sedailyCommonURL = "https://m.sedaily.com/RankAll"
const sedailyItemURL = "http://m.sedaily.com"

type Sedaily struct{}

func (c Sedaily) GetName() string {
	return "sedaily"
}

func (c Sedaily) GetGroup() string {
	return "1"
}

func (c Sedaily) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(sedailyCommonURL)
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

		aTag := sel.Find("a").First()

		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		title := aTag.Find("h2").Text()
		url := sedailyItemURL + href + "#_enliple"

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c Sedaily) GetContents(item *models.NewsItem) error {
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
	date := html.Find("div.article_info").Find("span.url_txt").Text()
	comp, _ := regexp.Compile("[0-9]{4}.[0-9]{2}.[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}")
	result := comp.FindString(date)

	if len(result) == 0 {
		item.Datetime = "today"
	} else {
		item.Datetime = result
	}

	wrapper := html.Find("div.article")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + " ")
		}
	})
	item.Contents = contents

	return nil
}
