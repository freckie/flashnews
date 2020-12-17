package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const etoday1202CommonURL = "https://www.etoday.co.kr/news/section/subsection?MID=1202"
const etoday1202ItemURL = "https://www.etoday.co.kr"

type Etoday1202 struct{}

func (c Etoday1202) GetName() string {
	return "etoday1202"
}

func (c Etoday1202) GetGroup() string {
	return "11"
}

func (c Etoday1202) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	httpReq, err := http.NewRequest("GET", etoday1202CommonURL, strings.NewReader(""))
	httpReq.Header.Add("Accept-Charset", "utf-8")
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
	wrapper := html.Find("div#list_W")
	items := wrapper.Find("div.cluster_text")
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
		title := aTag.Find("div.cluster_text_headline").Text()
		url := etoday1202ItemURL + href

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c Etoday1202) GetContents(item *models.NewsItem) error {
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

	item.Datetime = strings.Replace(html.Find("div.newsinfo > span").Text(), "입력", "", -1)

	// Parsing
	wrapper := html.Find("div.articleView > div")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		if idx > 0 {
			contents += utils.TrimAll(sel.Text())
		}
	})

	item.Contents = contents

	return nil
}
