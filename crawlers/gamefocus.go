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

const gameFocusCommonURL = "http://www.gamefocus.co.kr/html_file.php?file=normal_all_news.html"
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
	if number > 15 || number < 1 {
		_number = 15
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
	wrapper := html.Find("div.f_l").Find("table > tbody")
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

	item.Title = title
	item.Datetime = date
	item.Contents = contents

	return nil
}
