package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const bioSpectatorCommonURL = "http://www.biospectator.com/section/section_list.php?MID=10000"
const bioSpectatorItemURL = "http://www.biospectator.com"

type BioSpectator struct{}

func (c BioSpectator) GetName() string {
	return "biospectator"
}

func (c BioSpectator) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(bioSpectatorCommonURL)
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
	wrapper := html.Find("div.article_list")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		strong := sel.Find("strong.article_tit")
		aTag := strong.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := bioSpectatorItemURL + href

		date := sel.Find("span.date").Text()
		title := utils.TrimAll(aTag.Text())
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

func (c BioSpectator) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.article_view")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		if idx > 0 {
			contents += utils.TrimAll(sel.Text())
		}
	})

	item.Contents = contents
	return nil
}
