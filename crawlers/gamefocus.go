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

const gameFocusCommonURL = "http://www.gamefocus.co.kr/ajax_section.php?ajaxNum=2&ajaxLayer=section_ajax_layer_2&thread=22r02&pg=1&vals=%EC%9E%90%EB%8F%99%2C%EC%A0%84%EC%B2%B4%2C%EC%B4%9D1%2F15%EA%B0%9C%EC%B6%9C%EB%A0%A5%2C%EC%A0%9C%EB%AA%A954%EC%9E%90%EC%9E%90%EB%A6%84%2C%EB%B3%B8%EB%AC%B8270%EC%9E%90%EC%9E%90%EB%A6%84%2C%ED%88%AC%EB%AA%85%EC%83%89%2C%EB%88%84%EB%9D%BD0%EA%B0%9C%2C%EC%A0%84%EC%B2%B4%EB%89%B4%EC%8A%A4%EC%B6%9C%EB%A0%A5%2C%EC%9D%B4%EB%AF%B8%EC%A7%80%EA%B0%80%EB%A1%9C%ED%94%BD%EC%85%8080%2F53%2Csub_news_rows_21_1.html%2C%EC%9E%90%EB%8F%99%2C%ED%8E%98%EC%9D%B4%EC%A7%95%2C"
const gameFocusItemURL = "http://www.gamefocus.co.kr/"

var regexDate = regexp.MustCompile(`(년|월)`)
var regexTime = regexp.MustCompile(`(시)`)

type GameFocus struct{}

func (c GameFocus) GetName() string {
	return "gamefocus"
}

func (c GameFocus) GetGroup() string {
	return "5"
}

func (c GameFocus) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(gameFocusCommonURL)
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
	wrapper := html.Find("table > tbody")
	// wrapper := html.Find("div.f_l").Find("table > tbody")
	wrapper.ChildrenFiltered("tr").Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := gameFocusItemURL + href

		result[i] = models.NewsItem{
			Title:    "",
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c GameFocus) GetContents(item *models.NewsItem) error {
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
	title := utils.TrimAll(html.Find("div.detail_view").Find("h2").Text())
	date := utils.TrimAll(html.Find("span.font_11").First().Text())
	date = strings.Replace(regexDate.ReplaceAllString(date, `-`), "일", "", -1)
	date = strings.Replace(regexTime.ReplaceAllString(date, `:`), "분", "", -1)

	wrapper := html.Find("div#ct")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	if contents == "" {
		wrapper.Find("font").Each(func(idx int, sel *goquery.Selection) {
			contents += (utils.TrimAll(sel.Text()) + " ")
		})
	}

	item.Title = title
	item.Datetime = date
	item.Contents = contents

	return nil
}
