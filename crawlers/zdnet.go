package crawlers

import (
	"fmt"
	"net/http"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const zdnetCommonURL = "https://m.zdnet.co.kr/news_list.asp?zdknum=0000"
const zdnetItemURL = "https://m.zdnet.co.kr/"

type ZDNet struct{}

func (c ZDNet) GetName() string {
	return "zdnet"
}

func (c ZDNet) GetGroup() string {
	return "7"
}

func (c ZDNet) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	_req, err := http.NewRequest("GET", zdnetCommonURL, nil)
	if err != nil {
		return result, err
	}
	_req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36'")
	client := &http.Client{}
	req, err := client.Do(_req)
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
	wrapper := html.Find("div.sub_list")
	items := wrapper.Find("div.sub_li")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		div := sel.Find("div.txt")
		aTag := div.Find("a")
		href, ok := aTag.Attr("href")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		url := zdnetItemURL + href
		title := utils.TrimAll(aTag.Text())

		// date
		var date string
		sel.Find("p").Find("span").Each(func(i int, sel *goquery.Selection) {
			if i != 1 {
				return
			}
			date = sel.Text()
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

func (c ZDNet) GetContents(item *models.NewsItem) error {
	// Request
	_req, err := http.NewRequest("GET", item.URL, nil)
	if err != nil {
		return err
	}
	_req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36'")
	client := &http.Client{}
	req, err := client.Do(_req)
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
	wrapper := html.Find("div#content")
	contents := ""
	wrapper.Find("p").Each(func(idx int, sel *goquery.Selection) {
		temp := utils.TrimAll(sel.Text())
		temp, err := utils.ReadISO88591toUTF8(temp)
		if err != nil {
			return
		}
		contents += (temp + "")
	})

	item.Contents = contents
	return nil
}
