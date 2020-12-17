package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const theGuruCommonURL = "https://www.theguru.co.kr/news/article_list_all.html"
const theGuruItemURL = "https://www.theguru.co.kr"

type TheGuru struct{}

func (c TheGuru) GetName() string {
	return "TheGuru"
}

func (c TheGuru) GetGroup() string {
	return "11"
}

func (c TheGuru) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(theGuruCommonURL)
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
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := theGuruItemURL + href

		date := utils.TrimAll(aTag.Find("ul.ffd.art_info > li.date").Text())
		title := utils.TrimAll(aTag.Find("h2.clamp.c2").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c TheGuru) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#news_body_area")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		if idx > 0 {
			contents += utils.TrimAll(sel.Text())
		}
	})

	item.Contents = contents
	return nil
}
