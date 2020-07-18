package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const press9CommonURL = "http://www.press9.kr/news/articleList.html?page=1&sc_section_code=S1N12&sc_order_by=E"
const press9ItemURL = "http://www.press9.kr"

type Press9 struct{}

func (c Press9) GetName() string {
	return "Press9"
}

func (c Press9) GetGroup() string {
	return "10"
}

func (c Press9) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(press9CommonURL)
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
	wrapper := html.Find("section.article-list-content")
	items := wrapper.Find("div.table-row")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a.links")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := press9ItemURL + href

		date := sel.Find("div.list-dated").Text()
		date = strings.Split(date, " | ")[1]
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

func (c Press9) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#article-view-content-div")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "")
	})

	item.Contents = contents
	return nil
}
