package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const newsTownCommonURL = "http://www.newstown.co.kr/news/articleList.html?page=1&sc_section_code=S1N3"
const newsTownItemURL = "http://www.newstown.co.kr"

type NewsTown struct{}

func (c NewsTown) GetName() string {
	return "NewsTown"
}

func (c NewsTown) GetGroup() string {
	return "10"
}

func (c NewsTown) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(newsTownCommonURL)
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
	wrapper := html.Find("section.article-list-content")
	items := wrapper.Find("div.table-row")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a.links")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := newsTownItemURL + href

		date := sel.Find("div.list-dated").Text()
		date = strings.Split(date, " | ")[1]
		title := aTag.Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c NewsTown) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#_article")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "")
	})

	item.Contents = contents
	return nil
}
