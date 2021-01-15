package crawlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	urlPackage "net/url"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const bosaCommonURL = "http://m.bosa.co.kr/news/putNewsJson.php"
const bosaItemURL = "http://m.bosa.co.kr/news/articleView.html?idxno="

type bosaDataItem struct {
	IdxNo    string `json:"idxno"`
	Title    string `json:"title"`
	ViewDate string `json:"view_date"`
	ViewTime string `json:"view_time"`
}

type bosaData struct {
	News []bosaDataItem `json:"news"`
}

type Bosa struct{}

func (c Bosa) GetName() string {
	return "Bosa"
}

func (c Bosa) GetGroup() string {
	return "12"
}

func (c Bosa) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(bosaCommonURL)
	if err != nil {
		return result, err
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", req.StatusCode, req.Status)
	}

	// Decode JSON
	jsonData := bosaData{}
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
		url := bosaItemURL + item.IdxNo
		title, err := urlPackage.QueryUnescape(item.Title)
		if err != nil {
			title = ""
		}
		result[cnt] = models.NewsItem{
			Title:    strings.Replace(title, "+", " ", -1),
			URL:      url,
			Contents: "",
			Datetime: item.ViewDate + " " + item.ViewTime,
		}
		cnt++
	}
	return result, nil
}

func (c Bosa) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articleBody > div.body").Find("p")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "\n")
	})

	item.Contents = contents
	return nil
}
