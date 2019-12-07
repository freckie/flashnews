package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const infoStockDailyCommonURL = "http://www.infostockdaily.co.kr/news/articleList.html?view_type=sm"
const infoStockDailyItemURL = "http://www.infostockdaily.co.kr"

type InfoStockDaily struct{}

func (c InfoStockDaily) GetName() string {
	return "infostockdaily"
}

func (c InfoStockDaily) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(infoStockDailyCommonURL)
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
	items := wrapper.Find("div.list-block")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		div := sel.Find("div.list-titles")
		aTag := div.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := infoStockDailyItemURL + href

		date := sel.Find("div.list-dated").Text()
		date = strings.TrimSpace(strings.Split(date, "|")[2])
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

func (c InfoStockDaily) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#article-view-content-div")
	contents := ""
	pars := wrapper.Find("p")
	pars.Each(func(i int, sel *goquery.Selection) {
		contents += strings.TrimSpace(sel.Text())
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
