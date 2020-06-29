package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const newsPrimeYMHCommonURL = "http://m.newsprime.co.kr/section_list_writer.html?name=%EC%96%91%EB%AF%BC%ED%98%B8+%EA%B8%B0%EC%9E%90"
const newsPrimeYMHItemURL = "http://m.newsprime.co.kr/"

type NewsPrimeYMH struct{}

func (c NewsPrimeYMH) GetName() string {
	return "newsprimeymh"
}

func (c NewsPrimeYMH) GetGroup() string {
	return "9"
}

func (c NewsPrimeYMH) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 5 || number < 1 {
		_number = 5
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(newsPrimeYMHCommonURL)
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
	wrapper := html.Find("div.box01_0610_section")
	items := wrapper.Find("div.article_box_sl_section, .article_box_sl_section")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("div.title_text, .title_text_none").Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := newsPrimeYMHItemURL + href
		title := aTag.Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c NewsPrimeYMH) GetContents(item *models.NewsItem) error {
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
	hdTag := html.Find("div.hd")
	title := utils.TrimAll(hdTag.Find("p.tit").Text())
	date := utils.TrimAll(hdTag.Find("p.data").Text())

	wrapper := html.Find("div.stit2")
	contents := utils.TrimAll(wrapper.Text())

	item.Title = title
	item.Contents = contents
	item.Datetime = date

	return nil
}
