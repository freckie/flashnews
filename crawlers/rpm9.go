package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const rmp9CommonURL = "http://www.rpm9.com/news/section.html?id1=6"
const rmp9ItemURL = "http://www.rpm9.com"

type RPM9 struct{}

func (c RPM9) GetName() string {
	return "rpm9"
}

func (c RPM9) GetGroup() string {
	return "5"
}

func (c RPM9) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(rmp9CommonURL)
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
	wrapper := html.Find("ul.sub_newslist")
	items := wrapper.Find("li")
	cnt := 0
	items.Each(func(i int, sel *goquery.Selection) {
		if cnt >= _number {
			return
		}

		if class, ok := sel.Attr("class"); ok {
			if strings.Contains(class, "ad_text660") {
				return
			}
		}

		aTag := sel.Find("a.newstit")
		href, ok := aTag.Attr("href")
		if !ok {
			result[cnt] = models.NewsItem{}
			return
		}
		url := rmp9ItemURL + href

		title := aTag.Text()

		result[cnt] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}

		cnt++
	})

	return result, nil
}

func (c RPM9) GetContents(item *models.NewsItem) error {
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
	date := strings.Replace(html.Find("span.date").Text(), "발행일 : ", "", -1)

	wrapper := html.Find("div#articleBody")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "")
	})

	item.Contents = contents
	item.Datetime = date
	return nil
}
