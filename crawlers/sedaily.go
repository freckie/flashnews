package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const sedailyCommonURL = "https://m.sedaily.com/News/NewsAll"
const sedailyItemURL = "http://m.sedaily.com"

type Sedaily struct{}

func (c Sedaily) GetName() string {
	return "sedaily"
}

func (c Sedaily) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
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
	wrapper := html.Find("ul#newsList")
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
		title := aTag.Find("h2").Text()
		url := sedailyItemURL + href

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
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
	date := html.Find("div.article_info").Find("span").Text()
	item.Datetime = strings.Replace(date, "입력", "", -1)

	wrapper := html.Find("div.article")
	remove := wrapper.Find("script").Text()
	remove2 := wrapper.Find("div.ad_banner").Text()
	remove3 := wrapper.Find("span.sub_ad_banner4").Text()
	remove4 := wrapper.Find("div.al_cen").Text()
	contents := strings.Replace(wrapper.Text(), remove, "", -1)
	contents = strings.Replace(contents, remove2, "", -1)
	contents = strings.Replace(contents, remove3, "", -1)
	item.Contents = utils.TrimAll(strings.Replace(contents, remove4, "", -1))

	return nil
}
