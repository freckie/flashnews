package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const seoulWireCommonURL = "http://www.seoulwire.com/news/articleList.html?view_type=sm"
const seoulWireItemURL = "http://www.seoulwire.com"

type SeoulWire struct{}

func (c SeoulWire) GetName() string {
	return "seoulwire"
}

func (c SeoulWire) GetGroup() string {
	return "4"
}

func (c SeoulWire) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(seoulWireCommonURL)
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
	items := wrapper.Find("div.list-block")
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
		url := seoulWireItemURL + href

		dateSplit := strings.Split(sel.Find("div.list-dated").Text(), " | ")
		date := dateSplit[2]
		title := sel.Find("div.list-titles").Find("strong").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c SeoulWire) GetContents(item *models.NewsItem) error {
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
		contents += sel.Text()
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
