package crawlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const news1LatestCommonURL = "http://rest.news1.kr/archive/list?page=1&pg_per_cnt=10"
const news1LatestItemURL = "https://www.news1.kr/articles/?"

type news1LatestDataItem struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"pubdate1"`
}

type news1LatestData struct {
	Data []news1LatestDataItem `json:"data"`
}

type News1Latest struct{}

func (c News1Latest) GetName() string {
	return "news1latest"
}

func (c News1Latest) GetGroup() string {
	return "11"
}

func (c News1Latest) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(news1LatestCommonURL)
	if err != nil {
		return result, err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// Decode JSON
	jsonData := news1LatestData{}
	err = json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		return result, err
	}

	// Parsing
	cnt := 0
	for _, item := range jsonData.Data {
		if cnt >= number {
			break
		}
		url := news1LatestItemURL + item.ID
		result[cnt] = models.NewsItem{
			Title:    item.Title,
			URL:      url,
			Contents: item.Content,
			Datetime: item.Date,
		}
		cnt++
	}

	return result, nil
}

func (c News1Latest) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articles_detail")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + " ")
		}
	})

	item.Contents = utils.TrimAll(contents)

	return nil
}
