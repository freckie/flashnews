package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const ajuNewsCommonURL = "https://www.ajunews.com/util/sokbo_list?m=all&p=1"
const ajuNewsItemURL = "https://www.ajunews.com"

type AjuNews struct{}

func (c AjuNews) GetName() string {
	return "ajunews"
}

func (c AjuNews) GetGroup() string {
	return "8"
}

func (c AjuNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(ajuNewsCommonURL)
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
	wrapper := html.Find("div.sokbo_list").Find("ul")
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
		url := ajuNewsItemURL + href

		title := aTag.Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c AjuNews) GetContents(item *models.NewsItem) error {
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
	date := utils.TrimAll(strings.Split(html.Find("span.date").Text(), "입력 : ")[1])

	wrapper := html.Find("div#articleBody")

	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + " ")
		}
	})

	if len(contents) < 10 {
		wrapper = html.Find("div#articleBody > div")
		contents = ""
		wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
			if goquery.NodeName(sel) == "#text" {
				contents += (utils.TrimAll(sel.Text()) + " ")
			}
		})
	}

	item.Contents = contents
	item.Datetime = date
	return nil
}
