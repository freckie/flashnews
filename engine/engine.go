package engine

import (
	"flashnews/config"
	"flashnews/crawlers"
	"flashnews/utils"
	"log"
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

	return nil
}

// Run : Main Func
func (c *Engine) Run() {
	// var
	var isFirst map[string]bool
	var prevData map[string][]string
	var errorCount map[string]int
	for _, crawler := range c.Crawlers {
		name := crawler.GetName()
		isFirst[name] = true
		prevData[name] = nil
		errorCount[name] = 0
	}

	for _, crawler := range c.Crawlers {
		_name := crawler.GetName()

		// Parallel Crawling
		go func(name string) {

			// Main Loop in goroutine
			for {
				// Get New Items
				data, err := crawler.GetList(15)
				if err != nil {
					c.Logger.Printf("[ERROR] crawler.GetList() : crawler(%s) : %s", name, err)
					errorCount[name] += 1
					continue
				}

				// first time
				if isFirst[name] {
					prevData[name] = utils.MakeURLArray(data)
					time.Sleep(time.Millisecond * time.Duration(c.Cfg.Crawler.DelayTimer))
					continue
				}

				// Detect New Item
				for idx, _ := range data {
					if !utils.IsContain(data[idx].URL, prevData[name]) { // New Item
						if utils.TitleCond(data[idx]) { // if TitleCond true
							err = crawler.GetContents(&data[idx])
							if err != nil {
								c.Logger.Printf("[ERROR] crawler.GetContents() : crawler(%s) : idx(%d) : %s", name, idx, err)
								errorCount[name] += 1
								continue
							}

							detectKeywords, ok := utils.KeywordCond(data[idx], c.Cfg.Keywords)
							if ok && len(detectKeywords) >= 3 {
								go c.TG.SendMessage(data[idx], detectKeywords)
							}
						}
					}
				}

				time.Sleep(time.Millisecond * time.Duration(c.Cfg.Crawler.DelayTimer))
			}

		}(_name)
	}

}
