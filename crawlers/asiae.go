package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const asiaeCommonURL = "https://www.asiae.co.kr/realtime/sokbo_left.htm"
const asiaeItemURL = "https://www.asiae.co.kr"

type Asiae struct{}

func (c Asiae) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(asiaeCommonURL)
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
	wrapper := html.Find("div.ct.txtform")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a").First()

		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		title, ok := aTag.Attr("title")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		title, err = utils.ReadCP949(title)
		if err != nil {
			result[i] = models.NewsItem{}
			return
		}

		url := asiaeItemURL + href

		date := sel.Find("span.date").Text()

		result[i] = models.NewsItem{
			Title:    title,
			Keyword:  "",
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Asiae) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.txtbox")
	remove := wrapper.Find("table").Text()
	if remove != "" {
		remove, err = utils.ReadCP949(remove)
		if err != nil {
			return err
		}
	}

	contents, err := utils.ReadCP949(wrapper.Text())
	if err != nil {
		return err
	}

	contents = strings.Replace(contents, remove, "", -1)
	item.Contents = strings.TrimSpace(strings.Replace(contents, "\n", "", -1))

	return nil
}
