package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const newsisCommonURL = "http://www.newsis.com/realnews/"
const newsisItemURL = "http://www.newsis.com"

type Newsis struct{}

func (c Newsis) GetName() string {
	return "newsis"
}

func (c Newsis) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(newsisCommonURL)
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
	wrapper := html.Find("div.lst_p6.mgt21")
	items := wrapper.Find("li.p1_bundle")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		strong := sel.Find("strong.title")
		aTag := strong.Find("a")

		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		title := aTag.Text()
		url := newsisItemURL + href

		date := sel.Find("span.date").Text()
		date = strings.TrimSpace(strings.Split(date, " | ")[1])

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Newsis) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#textBody")
	remove := wrapper.Find("div.summary_view").Text()
	contents := strings.Replace(wrapper.Text(), remove, "", -1)

	removes := wrapper.Find("div.view_text")
	removes.Each(func(i int, sel *goquery.Selection) {
		contents = strings.Replace(contents, sel.Text(), "", -1)
	})

	item.Contents = utils.TrimAll(contents)

	return nil
}
