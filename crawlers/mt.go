package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mtCommonURL = "https://news.mt.co.kr/newsflash/newsflash.html?sec=all&listType=left"
const mtItemURL = "https://news.mt.co.kr/newsflash/frame_article.html?sec=mt&type=all&no="

type MT struct{}

func (c MT) GetName() string {
	return "mt"
}

func (c MT) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mtCommonURL)
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
	wrapper := html.Find("div.group")
	items := wrapper.Find("li.bundle")
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
		id := reForNumbers.FindString(href)
		title := strings.TrimSpace(aTag.Text())
		title, err = utils.ReadCP949(title)
		if err != nil {
			result[i] = models.NewsItem{}
			return
		}

		url := mtItemURL + id

		date := sel.Find("span.time").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c MT) GetContents(item *models.NewsItem) error {
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
	title := strings.TrimSpace(html.Find("div#article").Find("h1").Text())
	item.Title, err = utils.ReadCP949(title)
	if err != nil {
		return err
	}

	wrapper := html.Find("div#textBody")
	remove := wrapper.Find("table.article_photo").Text()

	if remove == "" {
		remove = ""
	} else {
		remove, err = utils.ReadCP949(remove)
		if err != nil {
			return err
		}
	}

	contents := wrapper.Text()
	contents, err = utils.ReadCP949(contents)
	if err != nil {
		return err
	}

	contents = strings.TrimSpace(strings.Replace(contents, remove, "", -1))
	item.Contents = strings.Replace(contents, "\n", "", -1)

	return nil
}
