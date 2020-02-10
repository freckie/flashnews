package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const sedailyGA05CommonURL = "https://www.sedaily.com/NewsList/GA05"
const sedailyGA05ItemURL = "http://www.sedaily.com"

type SedailyGA05 struct{}

func (c SedailyGA05) GetName() string {
	return "sedaily_ga05"
}

func (c SedailyGA05) GetGroup() string {
	return "5"
}

func (c SedailyGA05) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(sedailyGA05CommonURL)
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
	wrapper := html.Find("ul.news_list")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		var aTag *goquery.Selection
		if _, ok := sel.Attr("class"); ok {
			aTag = sel.Find("a")
		} else {
			aTag = sel.Find("div").Find("a")
		}

		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		title := aTag.Find("span").Text()
		url := sedailyGA05ItemURL + href
		date := sel.Find("span.letter").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c SedailyGA05) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.view_con")
	contents := ""
	wrapper.Contents().Each(func(idx int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + "")
		}
	})

	item.Contents = contents
	return nil
}
