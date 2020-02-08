package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dailyPharmCommonURL = "http://www.dailypharm.com/Users/News/NewsList.html"
const dailyPharmItemURL = "http://www.dailypharm.com/Users/News/"

type DailyPharm struct{}

func (c DailyPharm) GetName() string {
	return "dailypharm"
}

func (c DailyPharm) GetGroup() string {
	return "5"
}

func (c DailyPharm) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(dailyPharmCommonURL)
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
	wrapper := html.Find("div.seachResult").ChildrenFiltered("ul")
	items := wrapper.Find("li.newsList")
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
		url := dailyPharmItemURL + href

		var title, date string
		aTag.Find("div.listHead").Contents().Each(func(i int, sel2 *goquery.Selection) {
			if goquery.NodeName(sel2) == "#text" {
				title, err = utils.ReadISO88591(sel2.Text())
				if err != nil {
					title = ""
					return
				}
				return
			}

			if goquery.NodeName(sel2) == "span" {
				date = sel2.Text()
				return
			}
		})

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c DailyPharm) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.newsContents")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			text, err := utils.ReadISO88591(sel.Text())
			if err != nil {
				text = ""
				return
			}
			contents += (utils.TrimAll(text) + " ")
		}
	})

	item.Contents = contents
	return nil
}
