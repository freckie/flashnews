package crawlers

import (
	"errors"
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const dongaScienceCommonURL = "http://dongascience.donga.com/news.php?"
const dongaScienceItemURL = "http://dongascience.donga.com"

type DongaScience struct{}

func (c DongaScience) GetName() string {
	return "DongaScience"
}

func (c DongaScience) GetGroup() string {
	return "11"
}

func (c DongaScience) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Resolve URL up to 13 redirects.
	client := &http.Client{
		CheckRedirect: func() func(req *http.Request, via []*http.Request) error {
			redirects := 0
			return func(req *http.Request, via []*http.Request) error {
				if redirects > 100 {
					return errors.New("stopped after 12 redirects")
				}
				redirects++
				return nil
			}
		}(),
	}

	// Request
	req, err := http.NewRequest("GET", dongaScienceCommonURL, nil)
	if err != nil {
		return result, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Request.URL.String())
	panic("hihi")

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("[ERROR] Request Status Code : %d, %s", resp.StatusCode, resp.Status)
	}

	// HTML Load
	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return result, err
	}

	// Parsing
	items := html.Find("ul#newslist > li")
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
		url := dongaScienceItemURL + href

		date := sel.Find("span.date").Text()
		title := sel.Find("span.tit").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c DongaScience) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#article_body > p")
	contents := ""
	wrapper.Each(func(idx int, sel *goquery.Selection) {
		contents += (utils.TrimAll(sel.Text()) + "\n")
	})

	item.Contents = contents
	return nil
}
