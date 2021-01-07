package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const heraldCorpCommonURL = "http://biz.heraldcorp.com/list.php?ct=010106000000&ctm=19"
const heraldCorpItemURL = "http://biz.heraldcorp.com/"

type HeraldCorp struct{}

func (c HeraldCorp) GetName() string {
	return "HeraldCorp"
}

func (c HeraldCorp) GetGroup() string {
	return "12"
}

func (c HeraldCorp) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(heraldCorpCommonURL)
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

	// Header News
	header := html.Find("div.list_head")
	headerA := header.Find("a")
	headerTitle := utils.TrimAll(header.Find("div.list_head_title").Text())
	headerDate := utils.TrimAll(header.Find("div.list_date").Text())
	headerHref, ok := headerA.Attr("href")
	if !ok {
		result[0] = models.NewsItem{}
	}
	headerURL := heraldCorpItemURL + headerHref

	result[0] = models.NewsItem{
		Title:    headerTitle,
		URL:      headerURL,
		Contents: "",
		Datetime: headerDate,
	}

	// Parsing
	wrapper := html.Find("div.list").Find("ul")
	items := wrapper.Find("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number-1 {
			return
		}

		aTag := sel.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i+1] = models.NewsItem{}
			return
		}
		url := heraldCorpItemURL + href

		date := utils.TrimAll(sel.Find("div.l_date").Text())
		title := utils.TrimAll(aTag.Find("div.list_title").Text())

		result[i+1] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c HeraldCorp) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#articleText")
	contents := ""
	wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + " ")
	})

	item.Contents = contents
	return nil
}
