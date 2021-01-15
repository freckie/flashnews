package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const doctorsTimesCommonURL = "http://www.doctorstimes.com/news/articleList.html?view_type=sm"
const doctorsTimesItemURL = "http://www.doctorstimes.com"

type DoctorsTimes struct{}

func (c DoctorsTimes) GetName() string {
	return "DoctorsTimes"
}

func (c DoctorsTimes) GetGroup() string {
	return "12"
}

func (c DoctorsTimes) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(doctorsTimesCommonURL)
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

		div := sel.Find("div.list-titles")
		aTag := div.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := doctorsTimesItemURL + href

		date := utils.TrimAll(strings.Split(sel.Find("div.list-dated").Text(), "|")[2])
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

func (c DoctorsTimes) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#article-view-content-div > p")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "\n")
	})

	item.Contents = contents
	return nil
}
