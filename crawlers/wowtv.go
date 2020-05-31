package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const wowtvCommonURL = "http://www.wowtv.co.kr/NewsCenter/News/NewsList?subMenu=latest&menuSeq=459"
const wowtvItemURL = ""

type WowTV struct{}

func (c WowTV) GetName() string {
	return "wowtv"
}

func (c WowTV) GetGroup() string {
	return "9"
}

func (c WowTV) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(wowtvCommonURL)
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
	wrapper := html.Find("div.contain-list-news")
	items := wrapper.Find("div.article-news-list")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("div.contian-news").Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := wowtvItemURL + href

		date := utils.TrimAll(aTag.Find("span.date").Text())

		var title string
		aTag.Find("p.title-text").Contents().EachWithBreak(func(i int, sel *goquery.Selection) bool {
			if goquery.NodeName(sel) == "#text" {
				title = utils.TrimAll(sel.Text())
				return true
			}
			return false
		})
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

func (c WowTV) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#divNewsContent")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += utils.TrimAll(sel.Text() + " ")
		}
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
