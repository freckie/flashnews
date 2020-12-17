package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dNewsCommonURL = "https://www.dnews.co.kr/uhtml/autosec/D_S1N2_S2N20_1.html"
const dNewsItemURL = "https://www.dnews.co.kr"

type DNews struct{}

func (c DNews) GetName() string {
	return "DNews"
}

func (c DNews) GetGroup() string {
	return "11"
}

func (c DNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(dNewsCommonURL)
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
	items := html.Find("div.listBox_sub_main_s_l > ul > li")
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
		url := dNewsItemURL + href

		title := utils.TrimAll(aTag.Find("div.title").Text())
		date := utils.TrimAll(aTag.Find("div.news_date").Text())

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c DNews) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.newsCont > div.text > p")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		if idx > 0 {
			contents += utils.TrimAll(sel.Text())
		}
	})

	item.Contents = contents
	return nil
}
