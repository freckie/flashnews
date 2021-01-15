package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const medipharmHealthCommonURL = "http://www.medipharmhealth.co.kr/news/article_list_all.html"
const medipharmHealthItemURL = "http://www.medipharmhealth.co.kr"

type MedipharmHealth struct{}

func (c MedipharmHealth) GetName() string {
	return "MedipharmHealth"
}

func (c MedipharmHealth) GetGroup() string {
	return "12"
}

func (c MedipharmHealth) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(medipharmHealthCommonURL)
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
	items := html.Find("ul.art_list_all > li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a")
		title := utils.TrimAll(aTag.Find("h2.clamp.c2").Text())
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := medipharmHealthItemURL + href
		date := utils.TrimAll(aTag.Find("ul.art_info > li.date").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c MedipharmHealth) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#news_body_area > div")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	return nil
}
