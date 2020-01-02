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

const doc3CommonURL = "http://m.docdocdoc.co.kr/news/putNewsJson.php"
const doc3ItemURL = "http://m.docdocdoc.co.kr/news/articleView.html?idxno="

type doc3DataItem struct {
	IdxNo    string `json:"idxno"`
	Title    string `json:"title"`
	ViewData string `json:"view_date"`
}

type doc3Data struct {
	NewsCnt string         `json:"newsCnt"`
	News    []doc3DataItem `json:"news"`
}

type Doc3 struct{}

func (c Doc3) GetName() string {
	return "doc3"
}

func (c Doc3) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.PostForm(doc3CommonURL, url.Values{"page": {"1"}, "list_per_page": {"10"}})
	if err != nil {
		return result, err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// Decode JSON
	jsonData := doc3Data{}
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
		date := item.ViewData
		url := doc3ItemURL + item.IdxNo
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

func (c Doc3) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articleBody")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += sel.Text()
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents

	return nil
}
