package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mkVIP26CommonURL = "https://vip.mk.co.kr/newSt/news/news_list.php?sCode=26"
const mkVIP26ItemURL = "https:"

type MKVIP26 struct{}

func (c MKVIP26) GetName() string {
	return "mkvip26"
}

func (c MKVIP26) GetGroup() string {
	return "11"
}

func (c MKVIP26) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mkVIP26CommonURL)
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
		url := mkVIP26ItemURL + href

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

func (c MKVIP26) GetContents(item *models.NewsItem) error {
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
