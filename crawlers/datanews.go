package crawlers

import (
	"fmt"
	"net/http"
	"regexp"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dataNewsCommonURL = "https://m.datanews.co.kr/m/m_article_list_all.html"
const dataNewsItemURL = "https://m.datanews.co.kr/m/"

var dataNewsDatetimeRegex = regexp.MustCompile(`[0-9]{4}\.[0-9]{2}\.[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}`)

type DataNews struct{}

func (c DataNews) GetName() string {
	return "DataNews"
}

func (c DataNews) GetGroup() string {
	return "12"
}

func (c DataNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(dataNewsCommonURL)
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
	wrapper := html.Find("div.m01_araM1")
	items := wrapper.Find("dl.under_line")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("dd.m1 > strong > a")
		title := utils.TrimAll(aTag.Text())
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := dataNewsItemURL + href

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c DataNews) GetContents(item *models.NewsItem) error {
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

	// Date
	date := utils.TrimAll(html.Find("p.arvdate").Text())
	date = dataNewsDatetimeRegex.FindString(date)

	// Parsing
	wrapper := html.Find("div#news_body_area > div")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	item.Datetime = date
	return nil
}
