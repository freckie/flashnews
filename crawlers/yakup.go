package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const yakupCommonURL = "http://m.yakup.com"
const yakupItemURL = "http://m.yakup.com"

type Yakup struct{}

func (c Yakup) GetName() string {
	return "yakup"
}

func (c Yakup) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(yakupCommonURL)
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
	cnt := 0
	wrappers := html.Find("div.box-list-01")
	wrappers.Each(func(i int, sel *goquery.Selection) {
		if i >= 3 {
			return
		}

		items := sel.Find("li")
		items.Each(func(j int, sel2 *goquery.Selection) {
			if cnt >= _number {
				return
			}
			aTag := sel2.Find("a")
			href, ok := aTag.Attr("href")
			if !ok {
				result[cnt] = models.NewsItem{}
				return
			}
			url := yakupItemURL + href[1:]

			title := utils.TrimAll(aTag.Text())

			result[cnt] = models.NewsItem{
				Title:    title,
				URL:      url,
				Contents: "",
				Datetime: "",
			}
			cnt++
		})
	})

	return result, nil
}

func (c Yakup) GetContents(item *models.NewsItem) error {
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

	// Date
	dtTags := html.Find("dl.subdata").Find("dt")
	dtTags.Each(func(i int, sel *goquery.Selection) {
		if i == 2 {
			item.Datetime = sel.Text()
		}
	})

	// Parsing
	wrapper := html.Find("div.article_view")
	contents := utils.TrimAll(wrapper.Text())

	item.Contents = contents
	return nil
}
