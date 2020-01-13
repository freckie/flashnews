package crawlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const news1CommonURL = "http://www.news1.kr/ajax/ajax.php"
const news1ItemURL = "http://www.news1.kr/articles/?"

type news1DataItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Param3 string `json:"3"`
	Param6 string `json:"6"`
}

type news1Data struct {
	Msg  string          `json:"msg"`
	Data []news1DataItem `json:"data"`
}

type News1 struct{}

func (c News1) GetName() string {
	return "news1"
}

func (c News1) GetGroup() string {
	return "3"
}

func (c News1) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.PostForm(news1CommonURL, url.Values{
		"cmd":                 {"categories"},
		"op":                  {"categories_list"},
		"slimit":              {"1"},
		"elimit":              {"15"},
		"orderby":             {"pubdate_tsm"},
		"sort":                {"DESC"},
		"categories_sec":      {"parent"},
		"upper_categories_id": {"13"},
		"categories_id":       {"13"},
	})

	if err != nil {
		return result, err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// Decode JSON
	jsonData := news1Data{}
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
		title := item.Title
		date := item.Param6
		url := news1ItemURL + item.ID
		result[cnt] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: item.Param3,
			Datetime: date,
		}
		cnt++
	}

	return result, nil
}

func (c News1) GetContents(item *models.NewsItem) error {
	// HTML Load
	html, err := goquery.NewDocumentFromReader(strings.NewReader(item.Contents))
	if err != nil {
		return err
	}

	// Parsing
	contents := html.Text()
	remove := html.Find("table").Text()
	contents = strings.Replace(contents, remove, "", -1)

	item.Contents = utils.TrimAll(contents)

	return nil
}
