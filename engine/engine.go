package engine

import (
	"flashnews/config"
	"flashnews/crawlers"
	"flashnews/utils"
	"log"
	"sync"
	"time"
)

// Engine : Crawling Main Engine
type Engine struct {
	Logger   *log.Logger
	Cfg      *config.Config
	NewsCfg  *config.NewsConfig
	TG       *TGEngine
	Crawlers []crawlers.Crawler
}

// Init : Initialize Engine Driver
func (c *Engine) Init(logger *log.Logger, filePath string) error {
	var err error

	c.Logger = logger

	// Config
	c.Cfg, err = config.LoadConfig(filePath)
	if err != nil {
		return err
	}

	// News Config
	c.NewsCfg, err = config.LoadNewsConfig(c.Cfg.Crawler.InputPath2)
	if err != nil {
		return err
	}

	// Keywords
	c.Cfg.Keywords, err = config.LoadKeywords(c.Cfg.Crawler.InputPath)
	if err != nil {
		return err
	}

	// TG Engine
	c.TG = &TGEngine{}
	c.TG.Cfg = c.Cfg
	err = c.TG.GenerateBot()
	if err != nil {
		return err
	}

	// Setup Crawlers
	c.Crawlers = make([]crawlers.Crawler, 0)
	if c.NewsCfg.Asiae {
		c.Crawlers = append(c.Crawlers, crawlers.Asiae{})
	}
	if c.NewsCfg.Edaily {
		c.Crawlers = append(c.Crawlers, crawlers.Edaily{})
	}
	if c.NewsCfg.Etoday {
		c.Crawlers = append(c.Crawlers, crawlers.Etoday{})
	}
	if c.NewsCfg.MT {
		c.Crawlers = append(c.Crawlers, crawlers.MT{})
	}
	if c.NewsCfg.Sedaily {
		c.Crawlers = append(c.Crawlers, crawlers.Sedaily{})
	}

	c.Logger.Println("세팅 완료!")

	return nil
}

// Run : Main Func
func (c *Engine) Run() {
	// var
	isFirst := make(map[string]bool)
	prevData := make(map[string][]string)
	errorCount := make(map[string]int)
	for _, crawler := range c.Crawlers {
		name := crawler.GetName()
		isFirst[name] = true
		prevData[name] = nil
		errorCount[name] = 0
	}

	c.Logger.Println("메인 루프 시작.")

	var wait sync.WaitGroup
	defer wait.Wait()
	for _, _crawler := range c.Crawlers {
		_name := _crawler.GetName()

		// Concurrent Crawling
		wait.Add(1)
		go func(name string, crawler crawlers.Crawler) {
			defer wait.Done()

			c.Logger.Printf("%s 언론사 수집 시작.\n", name)

			// Main Loop in goroutine
			for {
				// Get New Items
				data, err := crawler.GetList(15)
				if err != nil {
					c.Logger.Printf("[ERROR] crawler.GetList() : crawler(%s) : %s", name, err)
					errorCount[name]++
					continue
				}

				// first time
				if isFirst[name] {
					prevData[name] = utils.MakeURLArray(data)
					isFirst[name] = false
					c.Logger.Printf("%s 언론사 첫 수집 완료.\n", name)
					time.Sleep(time.Millisecond * time.Duration(c.Cfg.Crawler.DelayTimer))
					continue
				}

				// Detect New Item
				for idx, _ := range data {
					if !utils.IsContain(data[idx].URL, prevData[name]) { // New Item

						c.Logger.Printf("[DEBUG] %s 새로운 item : %s\n", name, data[idx].Title)

						if utils.TitleCond(data[idx]) { // if TitleCond true
							c.Logger.Printf("[DEBUG] %s TitleCond 부합하는 item : %s\n", name, data[idx].Title)
							err = crawler.GetContents(&data[idx])
							if err != nil {
								c.Logger.Printf("[ERROR] crawler.GetContents() : crawler(%s) : idx(%d) : %s", name, idx, err)
								errorCount[name]++
								continue
							}

							detectKeywords, ok := utils.KeywordCond(data[idx], c.Cfg.Keywords)
							if ok && len(detectKeywords) >= 3 {
								c.Logger.Printf("[DEBUG] %s 조건에 부합하는 item : %s\n", name, data[idx].Title)
								go c.TG.SendMessage(data[idx], detectKeywords)
							}
						}
					}
				}

				prevData[name] = utils.MakeURLArray(data)
				c.Logger.Printf("[DEBUG] %s 한바퀴 완료.\n", name)
				time.Sleep(time.Millisecond * time.Duration(c.Cfg.Crawler.DelayTimer))
			}

		}(_name, _crawler) // end of goroutine
	}
}
