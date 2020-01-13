package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const ynaCommonURL = "https://www.yna.co.kr/news?site=navi_news"
const ynaItemURL = "https:"

type YNA struct{}

func (c YNA) GetName() string {
	return "yna"
}

func (c YNA) GetGroup() string {
	return "2"
}

func (c YNA) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(ynaCommonURL)
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
	wrapper := html.Find("div.headlines.headline-list")
	items := wrapper.Find("li.section02")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		dt := sel.Find("strong.news-tl")
		aTag := dt.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := ynaItemURL + href

		date := sel.Find("span.p-time").Text()
		title := aTag.Text()
		if err != nil {
			result[i] = models.NewsItem{}
			return
		}

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c YNA) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.article")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += sel.Text()
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
