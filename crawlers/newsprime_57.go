package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const newsPrime57CommonURL = "http://m.newsprime.co.kr/section_list.html?sec_no=57&menu=index"
const newsPrime57ItemURL = "http://m.newsprime.co.kr/"

type NewsPrime57 struct{}

func (c NewsPrime57) GetName() string {
	return "newsprime57"
}

func (c NewsPrime57) GetGroup() string {
	return "6"
}

func (c NewsPrime57) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(newsPrime57CommonURL)
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
	items := wrapper.Find("div.article_box_sl_section")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("div.title_text").Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := newsPrime57ItemURL + href
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

func (c NewsPrime57) GetContents(item *models.NewsItem) error {
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
