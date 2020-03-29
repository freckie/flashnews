package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const tfCommonURL = "http://news.tf.co.kr/list/all"
const tfItemURL = ""

type TF struct{}

func (c TF) GetName() string {
	return "tf"
}

func (c TF) GetGroup() string {
	return "8"
}

func (c TF) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 5 || number < 1 {
		_number = 5
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(tfCommonURL)
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
	wrapper := html.Find("div.txtList").Find("ul")
	items := wrapper.Find("li")
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
		url := tfItemURL + href

		date := utils.TrimAll(sel.Find("span").Text())

		result[i] = models.NewsItem{
			Title:    "",
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c TF) GetContents(item *models.NewsItem) error {
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
	title := utils.TrimAll(html.Find("div.articleTitle").Text())

	wrapper := html.Find("div.article")

	contents := ""
	wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	item.Title = title
	return nil
}
