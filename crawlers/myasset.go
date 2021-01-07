package crawlers

import (
	"fmt"
	"net/http"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const myAssetCommonURL = "https://www.myasset.com/myasset/research/rs_news/news_list.cmd?s_kind_code=00"
const myAssetItemURL = "https://www.myasset.com/myasset/research/rs_news/news_view.cmd?s_kind_code=%s&viewdate=%s"

type MyAsset struct{}

func (c MyAsset) GetName() string {
	return "MyAsset"
}

func (c MyAsset) GetGroup() string {
	return "12"
}

func (c MyAsset) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	req, err := http.Get(myAssetCommonURL)
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
	wrapper := html.Find("tbody#news_tboody")
	items := wrapper.Find("tr")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		tdTag := sel.Find("td.txtL")
		aTag := tdTag.Find("a")
		dataKind, ok := aTag.Attr("data-kind")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}
		dataSeq, ok := aTag.Attr("data-seq")
		if !ok {
			result[i] = models.NewsItem{}
			return
		}

		date := dataSeq
		title := utils.TrimAll(aTag.Contents().First().Text())
		url := fmt.Sprintf(myAssetItemURL, dataKind, strings.Replace(dataSeq, " ", "%20", -1))

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c MyAsset) GetContents(item *models.NewsItem) error {
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
	wrapper := html.Find("dd.view")
	contents := ""
	wrapper.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + " ")
		}
	})

	if length := len(utils.TrimAll(contents)); length == 0 {
		wrapper.Find("p").Each(func(i int, sel *goquery.Selection) {
			contents += (utils.TrimAll(sel.Text()) + " ")
		})
	}

	item.Contents = contents
	return nil
}
