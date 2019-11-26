package main

import (
	"fmt"

	_ "flashnews/crawlers"
	"flashnews/services"
	_ "flashnews/utils"
)

func main() {
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
