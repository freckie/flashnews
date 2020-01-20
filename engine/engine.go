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

	// Filters
	c.Cfg.Filters, err = config.LoadFilters(c.Cfg.Crawler.InputPath3)
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
	if c.NewsCfg.BizChosun {
		c.Crawlers = append(c.Crawlers, crawlers.BizChosun{})
	}
	if c.NewsCfg.FnNews {
		c.Crawlers = append(c.Crawlers, crawlers.FnNews{})
	}
	if c.NewsCfg.Hankyung {
		c.Crawlers = append(c.Crawlers, crawlers.Hankyung{})
	}
	if c.NewsCfg.InfoStockDaily {
		c.Crawlers = append(c.Crawlers, crawlers.InfoStockDaily{})
	}
	if c.NewsCfg.MK {
		c.Crawlers = append(c.Crawlers, crawlers.MK{})
	}
	if c.NewsCfg.MTN {
		c.Crawlers = append(c.Crawlers, crawlers.MTN{})
	}
	if c.NewsCfg.Newspim {
		c.Crawlers = append(c.Crawlers, crawlers.Newspim{})
	}
	if c.NewsCfg.YNA {
		c.Crawlers = append(c.Crawlers, crawlers.YNA{})
	}
	/* Group 3 */
	if c.NewsCfg.BioSpectator {
		c.Crawlers = append(c.Crawlers, crawlers.BioSpectator{})
	}
	if c.NewsCfg.DailyMedi {
		c.Crawlers = append(c.Crawlers, crawlers.DailyMedi{})
	}
	if c.NewsCfg.DocDocDoc {
		c.Crawlers = append(c.Crawlers, crawlers.Doc3{})
	}
	if c.NewsCfg.DoctorsNews {
		c.Crawlers = append(c.Crawlers, crawlers.DoctorsNews{})
	}
	if c.NewsCfg.MDToday {
		c.Crawlers = append(c.Crawlers, crawlers.MDToday{})
	}
	if c.NewsCfg.News1 {
		c.Crawlers = append(c.Crawlers, crawlers.News1{})
	}
	if c.NewsCfg.Newsis {
		c.Crawlers = append(c.Crawlers, crawlers.Newsis{})
	}
	if c.NewsCfg.NewsRun {
		c.Crawlers = append(c.Crawlers, crawlers.NewsRun{})
	}
	if c.NewsCfg.PaxnetNews {
		c.Crawlers = append(c.Crawlers, crawlers.PaxnetNews{})
	}
	if c.NewsCfg.Yakup {
		c.Crawlers = append(c.Crawlers, crawlers.Yakup{})
	}
	/* Group 4 */
	if c.NewsCfg.BusinessPost {
		c.Crawlers = append(c.Crawlers, crawlers.BusinessPost{})
	}
	if c.NewsCfg.DDaily {
		c.Crawlers = append(c.Crawlers, crawlers.DDaily{})
	}
	if c.NewsCfg.DT {
		c.Crawlers = append(c.Crawlers, crawlers.DT{})
	}
	if c.NewsCfg.GENews {
		c.Crawlers = append(c.Crawlers, crawlers.GENews{})
	}
	if c.NewsCfg.INews24 {
		c.Crawlers = append(c.Crawlers, crawlers.INews24{})
	}
	if c.NewsCfg.InTheNews {
		c.Crawlers = append(c.Crawlers, crawlers.InTheNews{})
	}
	if c.NewsCfg.Medipana {
		c.Crawlers = append(c.Crawlers, crawlers.Medipana{})
	}
	if c.NewsCfg.Newsway {
		c.Crawlers = append(c.Crawlers, crawlers.Newsway{})
	}
	if c.NewsCfg.Nspna {
		c.Crawlers = append(c.Crawlers, crawlers.Nspna{})
	}
	if c.NewsCfg.SeoulWire {
		c.Crawlers = append(c.Crawlers, crawlers.SeoulWire{})
	}
	if c.NewsCfg.TheBell {
		c.Crawlers = append(c.Crawlers, crawlers.TheBell{})
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
				if errorCount[name] >= 5 {
					c.Logger.Printf("[ERROR] crawler(%s) 에러 다발, 서버 문제로 추정, 당분간 대기.", name)
					time.Sleep(time.Millisecond * time.Duration(1000*60*5))
					errorCount[name] = 0
					continue
				}

				// Get New Items
				data, err := crawler.GetList(7)
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
						if utils.TitleCond(data[idx]) { // if TitleCond true
							err = crawler.GetContents(&data[idx])
							if err != nil {
								c.Logger.Printf("[ERROR] crawler.GetContents() : crawler(%s) : idx(%d) : %s", name, idx, err)
								errorCount[name]++
								continue
							}

							detectKeywords, ok := utils.KeywordCond(data[idx], c.Cfg.Keywords, c.Cfg.Filters)
							if ok && len(detectKeywords) >= c.Cfg.Crawler.KeywordDetectionNum {
								go c.TG.SendMessage(data[idx], detectKeywords)
							}
						}
					}
				}

				prevData[name] = utils.MakeURLArray(data)
				time.Sleep(time.Millisecond * time.Duration(c.Cfg.Crawler.DelayTimer))
			}

		}(_name, _crawler) // end of goroutine
	}
}
