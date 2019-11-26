package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"

	"github.com/PuerkitoBio/goquery"
)

const etodayCommonURL = "http://www.etoday.co.kr/main.php/news/flashnews/flash_list?MID=0&varPage=1"
const etodayItemURL = "http://www.etoday.co.kr/news/flashnews/flash_view?idxno="

type Etoday struct{}

func (c Etoday) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	httpReq, err := http.NewRequest("GET", etodayCommonURL, strings.NewReader(""))
	httpReq.Header.Add("Accept-Charset", "utf-8")
	req, err := http.DefaultClient.Do(httpReq)
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
	wrapper := html.Find("div.flash_tab_lst")
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
		id := reForNumbers.FindString(href)
		title := aTag.Text()
		url := etodayItemURL + id

		dateWrapper := sel.Find("div.flash_tab_press")
		date := dateWrapper.Find("span.flash_press").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Etoday) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div.articleView")
	remove := wrapper.Find("div.img_box_desc").Text()
	item.Contents = strings.TrimSpace(strings.Replace(wrapper.Text(), remove, "", -1))

	return nil
}
