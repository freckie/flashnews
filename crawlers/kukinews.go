package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const kukiNewsCommonURL = "http://m.kukinews.com/m/m_section.html?sec_no=66"
const kukiNewsItemURL = "http://m.kukinews.com"

type KukiNews struct{}

func (c KukiNews) GetName() string {
	return "kukinews"
}

func (c KukiNews) GetGroup() string {
	return "9"
}

func (c KukiNews) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number+1)

	// Request
	req, err := http.Get(kukiNewsCommonURL)
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

	// Headline
	headline := html.Find("div.headline")
	hlATag := headline.Find("a")
	hlTitle := utils.TrimAll(hlATag.Text())
	hlHref, ok := hlATag.Attr("href")
	if !ok {
		result[0] = models.NewsItem{}
	}
	hlURL := kukiNewsItemURL + hlHref
	result[0] = models.NewsItem{
		Title:    hlTitle,
		URL:      hlURL,
		Contents: "",
		Datetime: "",
	}

	// Parsing
	wrapper := html.Find("ul.lists")
	items := wrapper.ChildrenFiltered("li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i+1] = models.NewsItem{}
			return
		}
		url := kukiNewsItemURL + href

		title := utils.TrimAll(sel.Find("p.tit").Text())
		if err != nil {
			result[i+1] = models.NewsItem{}
			return
		}

		result[i+1] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: "",
		}
	})

	return result, nil
}

func (c KukiNews) GetContents(item *models.NewsItem) error {
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

	// Datetime
	date := utils.TrimAll(strings.Split(html.Find("div.byline").Text(), " | ")[1])
	item.Datetime = date

	// Parsing
	wrapper := html.Find("div#news_body_area")
	contents := ""
	wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
		contents += utils.TrimAll(sel.Text() + " ")
	})
	contents = utils.TrimAll(contents)

	item.Contents = contents
	return nil
}
