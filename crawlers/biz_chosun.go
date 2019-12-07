package crawlers

import (
	"flashnews/utils"
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"

	"github.com/PuerkitoBio/goquery"
)

const bizChosunCommonURL = "https://biz.chosun.com/svc/bulletin/index.html"
const bizChosunItemURL = "https://biz.chosun.com"

type BizChosun struct{}

func (c BizChosun) GetName() string {
	return "biz_chosun"
}

func (c BizChosun) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(bizChosunCommonURL)
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
	wrapper := html.Find("div.art_list_wrap")
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
		url := bizChosunItemURL + href

		date := sel.Find("span.time").Text()
		title := strings.Replace(aTag.Text(), date, "", -1)

		//title, err = utils.ReadCP949(title)
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

func (c BizChosun) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#news_body_id")

	title := html.Find("h1#news_title_text_id").Text()
	item.Title = title

	contents := ""
	pars := wrapper.Find("div.par")
	pars.Each(func(i int, sel *goquery.Selection) {
		contents += strings.TrimSpace(sel.Text())
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
