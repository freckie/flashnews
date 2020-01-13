package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mdtodayCommonURL = "http://www.mdtoday.co.kr/mdtoday/index.html"
const mdtodayItemURL = "http://m.mdtoday.co.kr"

type MDToday struct{}

func (c MDToday) GetName() string {
	return "mdtoday"
}

func (c MDToday) GetGroup() string {
	return "3"
}

func (c MDToday) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 14 || number < 1 {
		_number = 14
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mdtodayCommonURL)
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
	wrapper := html.Find("td#MainContent")
	items := wrapper.Find("div#box1")
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
		url := mdtodayItemURL + href

		date := sel.Find("font").Text()
		title, err := utils.ReadISO88591(aTag.Find("b").Text())
		if err != nil {
			fmt.Println("ISO 8895-1 ERROR :", err)
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

func (c MDToday) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articleBody")
	contents, _ := utils.ReadISO88591(wrapper.Text())

	item.Contents = contents
	return nil
}
