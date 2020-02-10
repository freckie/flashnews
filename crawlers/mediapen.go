package crawlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mediaPenCommonURL = "http://m.mediapen.com/news/lists/?menu=2&cate_cd1="
const mediaPenItemURL = "http://m.mediapen.com"

var regexMediaPenDate = regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}")

type MediaPen struct{}

func (c MediaPen) GetName() string {
	return "mediapen"
}

func (c MediaPen) GetGroup() string {
	return "5"
}

func (c MediaPen) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mediaPenCommonURL)
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
	wrapper := html.Find("ul.main-list")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		if class, ok := sel.Attr("class"); ok && strings.Contains(class, "ad_text660") {
			return
		}

		aTag := sel.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := mediaPenItemURL + href

		title := aTag.Text()
		title = regexMediaPenDate.ReplaceAllString(title, "")

		result[i] = models.NewsItem{
			Title:    utils.TrimAll(title),
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c MediaPen) GetContents(item *models.NewsItem) error {
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
	date := regexMediaPenDate.FindString(html.Find("a.news_author").Text())

	wrapper := html.Find("div.news_contents")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	item.Datetime = date

	return nil
}
