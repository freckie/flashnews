package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const medicalTimesCommonURL = "http://www.medicaltimes.com/Users/News/NewsList.html"
const medicalTimesItemURL = "http://www.medicaltimes.com"

type MedicalTimes struct{}

func (c MedicalTimes) GetName() string {
	return "medicaltimes"
}

func (c MedicalTimes) GetGroup() string {
	return "8"
}

func (c MedicalTimes) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(medicalTimesCommonURL)
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
	wrapper := html.Find("div.newsBody")
	items := wrapper.Find("div.NewsList")
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
		url := medicalTimesItemURL + href

		date := utils.TrimAll(sel.Find("span.nlRegDate").Text())

		result[i] = models.NewsItem{
			Title:    "",
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c MedicalTimes) GetContents(item *models.NewsItem) error {
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
	title, err := utils.ReadISO88591(html.Find("div.newsTitle_new").Text())
	if err != nil {
		title = ""
	}

	wrapper := html.Find("div.newsTxt")

	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			text, err := utils.ReadISO88591(sel.Text())
			if err != nil {
				text = ""
			}
			contents += (utils.TrimAll(text) + " ")
		}
	})

	item.Contents = contents
	item.Title = utils.TrimAll(title)
	return nil
}
