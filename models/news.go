package models

type NewsItem struct {
	Title    string `json:"title"`
	Keyword  string `json:"keyword"`
	URL      string `json:"url"` // url, primary key
	Contents string `json:"contents"`
	Datetime string `json:"datetime"`
}

type News struct {
	Company string     `json:"company"` // Company of News
	Items   []NewsItem `json:"items"`
}
