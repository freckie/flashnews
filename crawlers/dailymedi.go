package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dailyMediCommonURL = "http://www.dailymedi.com/ajax_section.php?ajaxNum=1&ajaxLayer=section_ajax_layer_1&thread=&numberpart=&file2=&file=normal_all_news.html&area=&pg=1&vals=%C0%DA%B5%BF%2C%C0%FC%C3%BC%2C%C3%D11%2F10%B0%B3%C3%E2%B7%C2%2C%C1%A6%B8%F142%C0%DA%C0%DA%B8%A7%2C%BA%BB%B9%AE250%C0%DA%C0%DA%B8%A7%2C%C5%F5%B8%ED%BB%F6%2C%B4%A9%B6%F40%B0%B3%2C%C0%FC%C3%BC%B4%BA%BD%BA%C3%E2%B7%C2%2C%C0%CC%B9%CC%C1%F6%B0%A1%B7%CE%C7%C8%BC%BF100%2Crows_photo_news07.html%2C%C0%DA%B5%BF%2C%C6%E4%C0%CC%C2%A1%2C&start_date=&end_date="
const dailyMediItemURL = "http://www.dailymedi.com/"

type DailyMedi struct{}

func (c DailyMedi) GetName() string {
	return "dailymedi"
}

func (c DailyMedi) GetGroup() string {
	return "3"
}

func (c DailyMedi) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(dailyMediCommonURL)
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
	items := wrapper.ChildrenFiltered("tr")
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
		url := dailyMediItemURL + href

		date := sel.Find("font").Text()

		result[i] = models.NewsItem{
			Title:    "",
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c DailyMedi) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#sub_center_contents2 > table > tbody")
	contents := ""
	title := ""
	wrapper.ChildrenFiltered("tr").Each(func(idx int, sel *goquery.Selection) {
		// Title (tr[3])
		if idx == 3 {
			sel.Find("tr").Each(func(idx2 int, sel2 *goquery.Selection) {
				if idx2 == 1 {
					title, err = utils.ReadCP949(replaceSpecialChar(sel2.Text()))
					if err != nil {
						title = ""
						return
					}
				}
			})
		} else if idx == 6 { // Contents (tr[6])
			table := sel.Find("table").Find("td#ct")
			contents, err = utils.ReadCP949(replaceSpecialChar(table.Text()))
			if err != nil {
				contents = ""
				return
			}
		} else {
			return
		}

	})

	item.Title = utils.TrimAll(strings.ReplaceAll(title, "혻", " "))
	item.Contents = utils.TrimAll(strings.ReplaceAll(contents, "혻", " "))
	return nil
}

func replaceSpecialChar(data string) string {
	result := strings.ReplaceAll(data, "“", "\"")
	result = strings.ReplaceAll(result, "”", "\"")
	result = strings.ReplaceAll(result, "‘", "'")
	result = strings.ReplaceAll(result, "’", "'")
	result = strings.ReplaceAll(result, "&nbsp;", " ")
	return result
}
