package main

import (
	"fmt"

	_ "flashnews/crawlers"
	"flashnews/services"
	_ "flashnews/utils"
)

func main() {
	/*
		var c crawlers.Crawler
		choice := "asiae"

		switch choice {
		case "edaily":
			c = crawlers.Edaily{}
		case "etoday":
			c = crawlers.Etoday{}
		case "sedaily":
			c = crawlers.Sedaily{}
		case "mt":
			c = crawlers.MT{}
		case "asiae":
			c = crawlers.Asiae{}
		}

		li, _ := crawlers.GetList(c, 15)
		for idx, _ := range li {
			err := crawlers.GetContents(c, &li[idx])
			if err != nil {
				fmt.Println(idx, "번 에러!", err)
			}
		}
		fmt.Println(li[3].Title)
		fmt.Println(li[3].Datetime)
		fmt.Println(li[3].Contents)
		fmt.Println(li[3].URL)

		keywords := []string{"대통령", "정부", "청와대", "환영만찬", "유명서점"}

	*/

	ce := services.CEngine{}
	err := ce.Init("config.json")
	if err != nil {
		fmt.Println("[INIT ERROR]", err)
	}
	fmt.Println(ce.Crawlers[0].GetList(15))

	/*
		for idx, item := range li {
			fmt.Println("=============", idx, "=============")
			fmt.Println("제목 :", item.Title)
			fmt.Println("내용 :", item.Contents)
			fmt.Println("TitleCond :", utils.TitleCond(item))
			fmt.Println(utils.KeywordCond(item, keywords))
		}*/
}
