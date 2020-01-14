package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const newswayCommonURL = "http://www.newsway.co.kr/news/lists"
const newswayItemURL = "http://www.newsway.co.kr"

type Newsway struct{}

func (c Newsway) GetName() string {
	return "newsway"
}

func (c Newsway) GetGroup() string {
	return "4"
}

func (c Newsway) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(newswayCommonURL)
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
	wrapper := html.Find("div.totalList")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		divTag := sel.Find("div.ritext")
		aTag := divTag.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := newswayItemURL + href

		date := divTag.Find("span").Text()
		title := sel.Find("strong").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Newsway) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articleText")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + " ")
		}
	})

	item.Contents = contents
	return nil
}
