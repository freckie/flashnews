package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const whosaengCommonURL = "http://www.whosaeng.com/sub_view.html?type=abs"
const whosaengItemURL = "http://www.whosaeng.com"

type Whosaeng struct{}

func (c Whosaeng) GetName() string {
	return "Whosaeng"
}

func (c Whosaeng) GetGroup() string {
	return "12"
}

func (c Whosaeng) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(whosaengCommonURL)
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
	wrapper := html.Find("div#news_list2_area")
	items := wrapper.Find("div.news_list2")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("dd.title > a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := whosaengItemURL + href

		title := utils.TrimAll(aTag.Text())
		date := utils.TrimAll(sel.Find("dd.write > span.wdate").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Whosaeng) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#textinput")
	contents := ""
	wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	return nil
}
