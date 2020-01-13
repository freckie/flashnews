package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dtCommonURL = "http://www.dt.co.kr/section.html?section_num=2900"
const dtItemURL = ""

type DT struct{}

func (c DT) GetName() string {
	return "dt"
}

func (c DT) GetGroup() string {
	return "4"
}

func (c DT) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(dtCommonURL)
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
	wrapper := html.Find("div.list_area")
	items := wrapper.Find("dl.article_list")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		dtTag := sel.Find("dt")
		aTag := dtTag.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := dtItemURL + href

		date, err := utils.ReadISO88591(sel.Find("span.date").Text())
		if err != nil {
			result[i] = models.NewsItem{}
			return
		} else {
			date = strings.Replace(date, "입력 ", "", -1)
		}

		title, err := utils.ReadISO88591(aTag.Text())
		if err != nil {
			result[i] = models.NewsItem{}
			return
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

func (c DT) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.art_txt")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			temp, err := utils.ReadISO88591(utils.TrimAll(sel.Text()))
			if err != nil {
				temp = ""
			}
			contents += (utils.TrimAll(temp) + " ")
		}
	})

	item.Contents = contents
	return nil
}
