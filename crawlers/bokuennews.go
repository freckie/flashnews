package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const bokuenNewsCommonURL = "http://www.bokuennews.com/news/article_list_all.html"
const bokuenNewsItemURL = "http://www.bokuennews.com/news/"

type BokuenNews struct{}

func (c BokuenNews) GetName() string {
	return "BokuenNews"
}

func (c BokuenNews) GetGroup() string {
	return "12"
}

func (c BokuenNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(bokuenNewsCommonURL)
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
	wrapper := html.Find("div.arl_008 > ul")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("p.title > a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := bokuenNewsItemURL + href

		title := utils.TrimAll(aTag.Text())
		date := utils.TrimAll(sel.Find("span.date").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c BokuenNews) GetContents(item *models.NewsItem) error {
	// Request
	httpReq, err := http.NewRequest("GET", item.URL, strings.NewReader(""))
	httpReq.Header.Add("Content-Type", "text/html; charset=utf-8;")
	httpReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req, err := http.DefaultClient.Do(httpReq)
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
	wrapper := html.Find("div#news_body_area")
	contents := ""
	wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	return nil
}
