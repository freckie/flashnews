package crawlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const edailyCommonURL = "https://www.edaily.co.kr/news/realtimenews"
const edailyItemURL = "https://www.edaily.co.kr/news/realtime/realtime_NewsRead.asp?newsid="

var regexNewsID = regexp.MustCompile(`[0-9]+`)

type Edaily struct{}

func (c Edaily) GetName() string {
	return "edaily"
}

func (c Edaily) GetGroup() string {
	return "1"
}

func (c Edaily) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	httpReq, err := http.NewRequest("GET", edailyCommonURL, strings.NewReader(""))
	httpReq.Header.Add("Content-Type", "text/html; charset=utf-8;")
	httpReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

	req, err := http.DefaultClient.Do(httpReq)
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
	wrapper := html.Find("div.news_list")
	items := wrapper.Find("dl")
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
		id := regexNewsID.FindString(href)
		url := edailyItemURL + id

		title := aTag.Find("span").Text()

		title, err = utils.ReadISO88591toUTF8(title)
		if err != nil {
			result[i] = models.NewsItem{}
		}

		result[i] = models.NewsItem{
			Title:    utils.TrimAll(title),
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c Edaily) GetContents(item *models.NewsItem) error {
	// div#viewcontent
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

	// Parse Datetime
	timeStr := html.Find("p.newsdate").Text()
	datetime := strings.Split(timeStr, " | ")[1]

	// Parsing
	wrapper := html.Find("div#viewcontent")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			text, err := utils.ReadCP949(sel.Text())
			if err != nil {
				text = ""
			}
			contents += (utils.TrimAll(text) + " ")
		}
	})

	item.Contents = contents
	item.Datetime = datetime

	return nil
}
