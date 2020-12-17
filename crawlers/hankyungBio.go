package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const hankyungBioCommonURL = "https://www.hankyung.com/bioinsight"
const hankyungBioItemURL = ""

type HankyungBio struct{}

func (c HankyungBio) GetName() string {
	return "HankyungBio"
}

func (c HankyungBio) GetGroup() string {
	return "11"
}

func (c HankyungBio) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number+3)

	// Request
	req, err := http.Get(hankyungBioCommonURL)
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

	// Parsing 1
	wrapper := html.Find("ul.news_top")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("h3.news_tit > a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := href

		title := utils.TrimAll(aTag.Text())
		contents := utils.TrimAll(sel.Find("p.lead").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: contents,
			Datetime: "",
		}
	})

	// Parsing 2
	wrapper = html.Find("div.news_list_wrap > ul.news_list")
	items = wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("div.txt_wrap > h3.news_tit > a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i+3] = models.NewsItem{}
			return
		}
		url := href

		title := utils.TrimAll(aTag.Text())
		contents := utils.TrimAll(sel.Find("div.txt_wrap > p.lead").Text())

		result[i+3] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: contents,
			Datetime: "",
		}
	})

	return result, nil
}

func (c HankyungBio) GetContents(item *models.NewsItem) error {
	return nil
}
