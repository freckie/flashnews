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

const newsRunCommonURL = "http://m.newsrun.co.kr/news/putNewsJson.php"
const newsRunItemURL = "http://m.newsrun.co.kr/news/articleView.html?idxno="

type newsRunDataItem struct {
	IdxNo    string `json:"idxno"`
	Title    string `json:"title"`
	ViewData string `json:"view_date"`
}

type newsRunData struct {
	NewsCnt string            `json:"newsCnt"`
	News    []newsRunDataItem `json:"news"`
}

type NewsRun struct{}

func (c NewsRun) GetName() string {
	return "newsrun"
}

func (c NewsRun) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.PostForm(newsRunCommonURL, url.Values{"page": {"1"}, "list_per_page": {"10"}})
	if err != nil {
		return result, err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// Decode JSON
	jsonData := newsRunData{}
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
		url := newsRunItemURL + item.IdxNo
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

func (c NewsRun) GetContents(item *models.NewsItem) error {
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
