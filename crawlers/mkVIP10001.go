package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mkVIP10001CommonURL = "https://vip.mk.co.kr/newSt/news/news_list.php?sCode=10001"
const mkVIP10001ItemURL = "https:"

type MKVIP10001 struct{}

func (c MKVIP10001) GetName() string {
	return "mkvip10001"
}

func (c MKVIP10001) GetGroup() string {
	return "11"
}

func (c MKVIP10001) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mkVIP10001CommonURL)
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
	wrapper := html.Find("table.table_6 > tbody")
	items := wrapper.Find("tr")
	cnt := 0
	items.Each(func(i int, sel *goquery.Selection) {
		if cnt >= _number {
			return
		}

		if i%5 == 0 {
			return
		}

		aTag := sel.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[cnt] = models.NewsItem{}
			cnt++
			return
		}
		url := mkVIP10001ItemURL + href

		date := sel.Find("td.t_11_brown").Text()
		title, err := utils.ReadCP949(aTag.Text())
		if err != nil {
			result[cnt] = models.NewsItem{}
			cnt++
			return
		}

		result[cnt] = models.NewsItem{
			Title:    utils.TrimAll(title),
			URL:      url,
			Contents: "",
			Datetime: date,
		}
		cnt++
	})

	return result, nil
}

func (c MKVIP10001) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#Conts")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			_contents, err := utils.ReadCP949(sel.Text())
			if err != nil {
				return
			}
			contents += (utils.TrimAll(_contents) + " ")
		}
	})

	item.Contents = contents
	return nil
}
