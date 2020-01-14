package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const medipanaCommonURL = "https://www.medipana.com/news/news_list_new.asp?MainKind=A&NewsKind=103&vCount=15&vKind=1"
const medipanaItemURL = "https://www.medipana.com"

type Medipana struct{}

func (c Medipana) GetName() string {
	return "medipana"
}

func (c Medipana) GetGroup() string {
	return "4"
}

func (c Medipana) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(medipanaCommonURL)
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
	wrapper := html.Find("div.totalNews")
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
		url := medipanaItemURL + strings.Replace(href, "../", "/", -1)

		dateSplit := strings.Split(sel.Find("span.infor").Text(), " | ")
		date := utils.TrimAll(dateSplit[1])
		title, err := utils.ReadCP949(sel.Find("span.tit").Text())
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

func (c Medipana) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.newsCon")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		temp, err := utils.ReadCP949(sel.Text())
		if err != nil {
			temp = ""
			return
		} else {
			temp = strings.Replace(temp, "혻", "", -1)
		}
		contents += (utils.TrimAll(temp) + " ")
	})

	if len(contents) == 0 {
		wrapper.Find("div").Each(func(idx int, sel *goquery.Selection) {
			temp, err := utils.ReadCP949(sel.Text())
			if err != nil {
				temp = ""
				return
			} else {
				temp = strings.Replace(temp, "혻", "", -1)
			}
			contents += (utils.TrimAll(temp) + " ")
		})
	}

	item.Contents = contents
	return nil
}
