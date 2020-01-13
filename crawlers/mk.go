package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mkCommonURL = "https://www.mk.co.kr/news/all/"
const mkItemURL = ""

type MK struct{}

func (c MK) GetName() string {
	return "mk"
}

func (c MK) GetGroup() string {
	return "2"
}

func (c MK) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mkCommonURL)
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
	wrapper := html.Find("div.list_area")
	items := wrapper.Find("dl.article_list")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		dt := sel.Find("dt.tit")
		aTag := dt.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := href

		date := sel.Find("span.date").Text()
		title, err := utils.ReadCP949(strings.TrimSpace(aTag.Text()))
		if err != nil {
			result[i] = models.NewsItem{}
			return
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

func (c MK) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.art_txt")
	contents, err := utils.ReadCP949(strings.TrimSpace(wrapper.Text()))
	if err != nil {
		return err
	}

	removes := wrapper.Find("div")
	removes.Each(func(idx int, sel *goquery.Selection) {
		remove, err := utils.ReadCP949(strings.TrimSpace(sel.Text()))
		if err != nil {
			remove = ""
		}
		contents = strings.Replace(contents, remove, "", -1)
	})
	removes2 := wrapper.Find("script")
	removes2.Each(func(idx int, sel *goquery.Selection) {
		remove2, err := utils.ReadCP949(strings.TrimSpace(sel.Text()))
		if err != nil {
			remove2 = ""
		}
		contents = strings.Replace(contents, remove2, "", -1)
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
