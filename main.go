package main

import (
	"fmt"

	"flashnews/crawlers"
	"flashnews/utils"
)

func main() {
	var c crawlers.Crawler
	choice := "mt"

	switch choice {
	case "edaily":
		c = crawlers.Edaily{}
	case "etoday":
		c = crawlers.Etoday{}
	case "seoul":
		c = crawlers.Seoul{}
	case "mt":
		c = crawlers.MT{}
	}

	li, _ := crawlers.GetList(c, 15)
	for idx, _ := range li {
		crawlers.GetContents(c, &li[idx])
	}
	/*
		fmt.Println(li[3].Title)
		fmt.Println(li[3].Datetime)
		fmt.Println(li[3].Contents)
		fmt.Println(li[3].URL)
	*/

	keywords := []string{"대통령", "정부", "청와대", "환영만찬", "유명서점"}

	for idx, item := range li {
		fmt.Println("=============", idx, "=============")
		fmt.Println("제목 :", item.Title)
		fmt.Println("내용 :", item.Contents)
		fmt.Println("TitleCond :", utils.TitleCond(item))
		fmt.Println(utils.KeywordCond(item, keywords))
	}
}
