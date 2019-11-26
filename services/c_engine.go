package services

import (
	"flashnews/config"
	"flashnews/crawlers"
)

// CEngine : Crawling Engine
type CEngine struct {
	Cfg      *config.Config
	NewsCfg  *config.NewsConfig
	TG       *TGEngine
	Crawlers []crawlers.Crawler
}

// Init : Initialize Engine Driver
func (c *CEngine) Init(filePath string) error {
	var err error

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
