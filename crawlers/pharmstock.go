package crawlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const pharmStockCommonURL = "http://m.pharmstock.co.kr/news/putNewsJson.php"
const pharmStockItemURL = "http://m.pharmstock.co.kr/news/articleView.html?idxno="

type pharmStockData struct {
	News []pharmStockDataItem `json:"news"`
}

type pharmStockDataItem struct {
	IdxNo    string `json:"idxno"`
	Title    string `json:"title"`
	ViewDate string `json:"view_date"`
	ViewTime string `json:"view_time"`
}

type PharmStock struct{}

func (c PharmStock) GetName() string {
	return "PharmStock"
}

func (c PharmStock) GetGroup() string {
	return "10"
}

func (c PharmStock) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.PostForm(pharmStockCommonURL, url.Values{"page": {"1"}, "sc_order_by": {"E"}})
	if err != nil {
		return result, err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// Decode JSON
	jsonData := pharmStockData{}
	err = json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		return result, err
	}

	// Parsing
	cnt := 0
	for _, item := range jsonData.News {
		if cnt >= number {
			break
		}
		title, err := url.QueryUnescape(item.Title)
		if err != nil {
			title = ""
		}
		date := item.ViewDate + " " + item.ViewTime
		url := pharmStockItemURL + item.IdxNo
		result[cnt] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
		cnt++
	}

	return result, nil
}

func (c PharmStock) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.body.word_break")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "")
	})

	item.Contents = contents
	return nil
}
