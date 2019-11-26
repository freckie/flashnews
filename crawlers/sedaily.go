package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"

	"github.com/PuerkitoBio/goquery"
)

const sedailyCommonURL = "https://m.sedaily.com/News/NewsAll"
const sedailyItemURL = "http://m.sedaily.com"

type Sedaily struct{}

func (c Sedaily) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(sedailyCommonURL)
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

		dlTag := sel.Find("dl")
		aTag := dlTag.Find("a").First()

		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		title := aTag.Text()
		url := sedailyItemURL + href

		date := dlTag.Find("span.letter").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Sedaily) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.view_con.first_view_con")
	remove := wrapper.Find("table").Text()
	remove2 := wrapper.Find("div.sub_ad_banner10_m").Text()
	remove3 := wrapper.Find("span").Text()
	contents := strings.Replace(wrapper.Text(), remove, "", -1)
	contents = strings.Replace(contents, remove2, "", -1)
	item.Contents = strings.TrimSpace(strings.Replace(contents, remove3, "", -1))

	return nil
}
