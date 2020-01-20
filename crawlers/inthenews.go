package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const inTheNewsCommonURL = "http://inthenews.co.kr/latest-news/"
const inTheNewsItemURL = ""

type InTheNews struct{}

func (c InTheNews) GetName() string {
	return "inthenews"
}

func (c InTheNews) GetGroup() string {
	return "4"
}

func (c InTheNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(inTheNewsCommonURL)
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
	wrapper := html.Find("div.pt-cv-wrapper")
	items := wrapper.Find("div.pt-cv-ifield")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		h4Tag := sel.Find("h4.pt-cv-title")
		aTag := h4Tag.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := inTheNewsItemURL + href

		date := ""
		title := aTag.Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c InTheNews) GetContents(item *models.NewsItem) error {
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
	date := utils.TrimAll(html.Find("div.post-date").Text())
	item.Datetime = date

	wrapper := html.Find("div.pf-content")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += sel.Text()
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
