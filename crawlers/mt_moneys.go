package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const mtMoneySCommonURL = "http://moneys.mt.co.kr/news/mwList.php?code=w0000&code2=w0100"
const mtMoneySItemURL = "http://mnb.moneys.mt.co.kr/mnbview.php"

type MTMoneyS struct{}

func (c MTMoneyS) GetName() string {
	return "mtmoneys"
}

func (c MTMoneyS) GetGroup() string {
	return "5"
}

func (c MTMoneyS) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 15 || number < 1 {
		_number = 15
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(mtMoneySCommonURL)
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
	wrapper := html.Find("ul.group.mgt13")
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
		title := strings.TrimSpace(sel.Find("strong.tit").Text())
		title, err = utils.ReadCP949(title)
		if err != nil {
			result[i] = models.NewsItem{}
			return
		}

		url := strings.Replace(href, "http://moneys.mt.co.kr/news/mwView.php", mtMoneySItemURL, -1)

		date := sel.Find("span.date").Text()

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c MTMoneyS) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("div#textBody") // > div")
	contents := ""
	// wrapper.Contents().Each(func(idx int, sel *goquery.Selection) {
	// 	if nn := goquery.NodeName(sel); nn == "#text" || nn == "p" {
	// 		text, err := utils.ReadCP949(sel.Text())
	// 		if err != nil {
	// 			return
	// 		}
	// 		contents += (utils.TrimAll(text) + " ")
	// 	}
	// })
	contents, err = utils.ReadCP949(wrapper.Text())
	if err != nil {
		contents = ""
	}
	contents = utils.TrimAll(contents)

	item.Contents = contents

	return nil
}
