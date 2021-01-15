package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const hitNewsCommonURL = "http://www.hitnews.co.kr/"
const hitNewsItemURL = "http://www.hitnews.co.kr"

type HITNews struct{}

func (c HITNews) GetName() string {
	return "HITNews"
}

func (c HITNews) GetGroup() string {
	return "12"
}

func (c HITNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(hitNewsCommonURL)
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
	wrapper := html.Find("section.content > div#skin-13")
	items := wrapper.Find("div.item")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a.auto-titles")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := hitNewsItemURL + href

		title := utils.TrimAll(aTag.Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c HITNews) GetContents(item *models.NewsItem) error {
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

	// Date parsing
	date := ""
	html.Find("ul.infomation > li").Each(func(i int, sel *goquery.Selection) {
		if i == 1 {
			date = utils.TrimAll(strings.Replace(sel.Text(), "입력", "", -1))
		}
	})

	// Parsing
	wrapper := html.Find("article#article-view-content-div")
	contents := ""
	wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	item.Datetime = date
	return nil
}
