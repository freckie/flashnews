package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dealSiteCommonURL = "https://dealsite.co.kr/categories?page=1"
const dealSiteItemURL = "https://dealsite.co.kr"

type DealSite struct{}

func (c DealSite) GetName() string {
	return "DealSite"
}

func (c DealSite) GetGroup() string {
	return "10"
}

func (c DealSite) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(dealSiteCommonURL)
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
	wrapper := html.Find("div.list")
	items := wrapper.Find("div.list-article")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a.dyn.std")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := dealSiteItemURL + href

		date := utils.TrimAll(sel.Find("div.pubdate").Text())
		title := aTag.Text()
		contents := sel.Find("span.sneakpeek").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: contents,
			Datetime: date,
		}
	})

	return result, nil
}

func (c DealSite) GetContents(item *models.NewsItem) error {
	return nil
}
