package crawlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"flashnews/models"
	"flashnews/utils"

	"github.com/PuerkitoBio/goquery"
)

const kiwoomCommonURL = "https://www2.kiwoom.com/nkw.TD6000News.do"
const kiwoomItemURL = "https://www2.kiwoom.com/nkw.TD6000NewsCont.do"

type Kiwoom struct{}

func (c Kiwoom) GetName() string {
	return "kiwoom"
}

func (c Kiwoom) GetGroup() string {
	return "11"
}

func (c Kiwoom) GetList(number int) ([]models.NewsItem, error) {
	// Number
	var _number int
	if number > 10 || number < 1 {
		_number = 10
	} else {
		_number = number
	}
	result := make([]models.NewsItem, _number)

	// Request
	httpReq, err := http.NewRequest("GET", kiwoomCommonURL, strings.NewReader(""))
	httpReq.Header.Add("Content-Type", "text/html; charset=utf-8;")
	httpReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

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
	newListArea := html.Find("div.newListArea")

	// Parsing
	wrapper := html.Find("table#oTable > tbody")
	items := wrapper.ChildrenFiltered("tr")
	items.Each(func(i int, sel *goquery.Selection) {
		if i >= _number {
			return
		}

		aTag := sel.Find("td.ldata > a")
		title := utils.TrimAll(aTag.Text())

		divIDs := []string{
			fmt.Sprintf("gubn_%d", i), // for supplier
			fmt.Sprintf("subg_%d", i), // for supplier
			fmt.Sprintf("date_%d", i), // date
			fmt.Sprintf("time_%d", i), // time
			fmt.Sprintf("seqn_%d", i), // for hcode
		}
		supplier := utils.TrimAll(newListArea.Find("div#"+divIDs[0]).Text()) +
			utils.TrimAll(newListArea.Find("div#"+divIDs[1]).Text())
		date := utils.TrimAll(newListArea.Find("div#"+divIDs[2]).Text()) + " " +
			utils.TrimAll(newListArea.Find("div#"+divIDs[3]).Text())
		seqn := utils.TrimAll(newListArea.Find("div#" + divIDs[4]).Text())

		hCode := date + seqn
		for _, tok := range []string{":", "/", " "} {
			hCode = strings.ReplaceAll(hCode, tok, "")
		}

		// Warning! this is not the url actually.
		url := supplier + "/" + hCode

		result[i] = models.NewsItem{
			Title:    title,
			URL:      url,
			Contents: "",
			Datetime: date,
		}
	})

	return result, nil
}

func (c Kiwoom) GetContents(item *models.NewsItem) error {
	// Request
	_formData := strings.Split(item.URL, "/")
	req, err := http.PostForm(kiwoomItemURL, url.Values{
		"supplier": {_formData[0]},
		"hcode":    {_formData[1]},
	})
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
	wrapper := html.Find("body")

	// Parsing
	contents := ""
	wrapper.Contents().Each(func(idx int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "#text" {
			contents += (utils.TrimAll(sel.Text()) + "")
		}
	})
	// contents := utils.TrimAll(html.Text())

	item.Contents = contents
	item.URL = ""
	return nil
}
