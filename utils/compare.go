package utils

import (
	"flashnews/models"
)

func MakeURLArray(items []models.NewsItem) []string {
	result := make([]string, len(items))

	for idx, item := range items {
		result[idx] = item.URL
	}

	return result
}

func IsContain(target string, arr []string) bool {
	for idx, _ := range arr {
		if arr[idx] == target {
			return true
		}
	}
	return false
}
