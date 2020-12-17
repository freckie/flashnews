package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const econovillCommonURL = "https://www.econovill.com/news/articleList.html?box_idxno=20&view_type=sm"
const econovillItemURL = "https://www.econovill.com"

type Econovill struct{}

func (c Econovill) GetName() string {
	return "Econovill"
}

func (c Econovill) GetGroup() string {
	return "11"
}

func (c Econovill) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(econovillCommonURL)
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
	items := html.Find("section#section-list > ul > li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("h4.titles > a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := econovillItemURL + href

		title := utils.TrimAll(aTag.Text())

		date := ""
		sel.Find("span.byline > em").Each(func(j int, sel2 *goquery.Selection) {
			if j == 2 {
				date = utils.TrimAll(sel2.Text())
			}
		})

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Econovill) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("article#article-view-content-div > p")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		if idx > 0 {
			contents += utils.TrimAll(sel.Text())
		}
	})

	item.Contents = contents
	return nil
}
