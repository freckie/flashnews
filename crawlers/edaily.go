package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const edailyCommonURL = "https://www.edaily.co.kr/news/realtime/realtime_NewsList_1.asp"
const edailyItemURL = "https://www.edaily.co.kr/news/realtime/realtime_NewsRead.asp?newsid="

type Edaily struct{}

func (c Edaily) GetName() string {
	return "edaily"
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
	wrapper := html.Find("ul")
	items := wrapper.Find("li")
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
		id := reForNumbers.FindString(href)
		url := edailyItemURL + id

		title := aTag.AttrOr("title", aTag.Text())
		date := sel.Find("span").Text()

		title, err = utils.ReadCP949(title)
		if err != nil {
			result[i] = models.NewsItem{}
		}

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
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

	// Parsing
	wrapper := html.Find("div#viewcontent")
	remove := wrapper.Find("div.ovh").Text()
	if remove != "" {
		remove, err = utils.ReadCP949(remove)
		if err != nil {
			return err
		}
	}
	remove2 := wrapper.Find("font").Text()
	if remove2 != "" {
		remove2, err = utils.ReadCP949(remove2)
		if err != nil {
			return err
		}
	}

	contents := wrapper.Text()
	contents, err = utils.ReadCP949(contents)
	if err != nil {
		return err
	}

	contents = strings.Replace(contents, remove, "", -1)
	contents = strings.Replace(contents, remove2, "", -1)
	contents = strings.Replace(contents, "\n", "", -1)
	contents = strings.Replace(contents, "\t", "", -1)
	item.Contents = strings.TrimSpace(strings.Replace(contents, "  ", "", -1))
	//ovh mt20
	return nil
}
