package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const healthInNewsCommonURL = "http://www.healthinnews.co.kr/news/articleList.html?view_type=sm"
const healthInNewsItemURL = "http://www.healthinnews.co.kr"

type HealthInNews struct{}

func (c HealthInNews) GetName() string {
	return "HealthInNews"
}

func (c HealthInNews) GetGroup() string {
	return "12"
}

func (c HealthInNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(healthInNewsCommonURL)
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
	wrapper := html.Find("section#section-list > ul.type2")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("h4.titles > a")
		title := utils.TrimAll(aTag.Text())
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := healthInNewsItemURL + href

		// date := utils.TrimAll(strings.Split(sel.Find("span.byline").Text(), " ")[2])

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c HealthInNews) GetContents(item *models.NewsItem) error {
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

	date := utils.TrimAll(html.Find("div.info-group").Find("i.icon-clock-o").Parent().Text())
	date = strings.Replace(date, "입력 ", "", -1)

	// Parsing
	wrapper := html.Find("article#article-view-content-div > p")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "\n")
	})

	item.Contents = contents
	item.Datetime = date
	return nil
}
