package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const nspnaCommonURL = "http://www.nspna.com/news/?cid=20"
const nspnaItemURL = "http://www.nspna.com"

type Nspna struct{}

func (c Nspna) GetName() string {
	return "nspna"
}

func (c Nspna) GetGroup() string {
	return "4"
}

func (c Nspna) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(nspnaCommonURL)
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
	wrapper := html.Find("div#news_list")
	items := wrapper.Find("div.news_panel")
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
		url := nspnaItemURL + href

		dateSplit := strings.Split(sel.Find("div.info").Text(), " | ")
		date := dateSplit[1]
		title := sel.Find("div.subject").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Nspna) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#CmAdContent")
	contents := ""
	wrapper.Find("p.section_txt").Each(func(idx int, sel *goquery.Selection) {
		contents += sel.Text()
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
