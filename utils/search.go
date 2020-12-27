package utils

import (
	"regexp"
	"strings"

	"flashnews/models"
)

// TitleCond : 제목 조건에 부합하는지.
func TitleCond(item models.NewsItem) bool {
	var cond1, cond2, cond3 bool
	title := item.Title

	// Cond1 : 첫 단어 이후 ', ", \''로 시작하는지.
	cond1, _ = regexp.MatchString("^([0-9a-zA-Z가-힣]+ ?[,，\"'“˙])", title)

	// Cond2 : 처음이 '[B [S [단 [특 [종' 와 같은지
	// Cond3 : 처음이 '?'와 같은지
	if temp := strings.TrimSpace(title); len(temp) > 0 {
		substr := StringSplit(title, 2)
		cond2 = ((substr == "[B") || (substr == "[S") || (substr == "[단") || (substr == "[특") || (substr == "[종"))
		substr2 := StringSplit(title, 1)
		cond3 = (substr2 == "?")
	} else {
		return false
	}

	return cond1 || cond2 || cond3
}

// KeywordCond : 제목과 본문에 키워드가 있는지
func KeywordCond(item models.NewsItem, keywords []string, filters []string) ([]string, bool) {
	match := make([]string, 0)
	isMatch := false

	_title := strings.Replace(item.Title, " ", "", -1)
	_contents := strings.Replace(item.Contents, " ", "", -1)
	_firstWord := strings.Split(_title, " ")[0]

	// Check Filters
	for _, filter := range filters {
		if strings.Contains(_firstWord, filter) {
			return match, false
		}
	}

	// Check Keywords
	for _, keyword := range keywords {
		if strings.Contains(_title, keyword) || strings.Contains(_contents, keyword) {
			match = append(match, keyword)
			isMatch = true
		}
	}

	return match, isMatch
}

// GetDateString : 해당 문자열에서 날짜 문자열을 추출
func GetDateString(data string) (string, bool) {
	comp, _ := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}")
	result := comp.FindString(data)

	if len(result) == 0 {
		return "", false
	} else {
		return result, true
	}
}
