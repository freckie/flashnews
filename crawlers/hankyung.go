package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const hankyungCommonURL = "https://www.hankyung.com/all-news"
const hankyungItemURL = ""

type Hankyung struct{}

func (c Hankyung) GetName() string {
	return "hankyung"
}

func (c Hankyung) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(hankyungCommonURL)
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
	wrapper := html.Find("ul.article_list")
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
		url := href

		date := sel.Find("p.time").Text()
		title := aTag.Text()

		if err != nil {
			result[i] = models.NewsItem{}
		}

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Hankyung) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articletxt")
	contents := wrapper.Text()

	removes := wrapper.Find("div.wrap_img")
	removes.Each(func(idx int, sel *goquery.Selection) {
		contents = strings.Replace(contents, sel.Text(), "", -1)
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
