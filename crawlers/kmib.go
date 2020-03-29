package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const kmibCommonURL = "http://news.kmib.co.kr/article/list_all.asp?sid1=all"
const kmibItemURL = "http://news.kmib.co.kr/article/"

type KMIB struct{}

func (c KMIB) GetName() string {
	return "kmib"
}

func (c KMIB) GetGroup() string {
	return "8"
}

func (c KMIB) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(kmibCommonURL)
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
	wrapper := html.Find("div.nws_list_all")
	items := wrapper.Find("dl.nws")
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
		url := kmibItemURL + href

		title, err := utils.ReadISO88591(aTag.Text())
		if err != nil {
			title = ""
		}
		date := utils.TrimAll(sel.Find("dd.date").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c KMIB) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articleBody")

	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			text, err := utils.ReadISO88591(sel.Text())
			if err != nil {
				text = ""
			}
			contents += (utils.TrimAll(text) + " ")
		}
	})

	item.Contents = contents
	return nil
}
