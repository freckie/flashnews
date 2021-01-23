package engine

import (
	"flashnews/config"
	"flashnews/crawlers"
	"flashnews/utils"
	"log"
	"strings"
	"sync"
	"time"
)

var CondExceptionCrawlers = []string{"gamefocus", "gamefocus22r09"}
var CondOnlyContentsCrawlers = []string{} // Ignores TitleCond

// Engine : Crawling Main Engine
type Engine struct {
	Logger   *log.Logger
	Cfg      *config.Config
	NewsCfg  *config.NewsConfig
	TG       *TGEngine
	SE       *SoundEngine
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
	c.Logger.Println("Telegram Engine 세팅 완료!")

	// Sound Engine
	c.SE = &SoundEngine{}
	err = c.SE.Init(c.Cfg.Sound.On, c.Cfg.Sound.FilePath)
	if err != nil {
		return err
	}
	c.Logger.Println("Sound Engine 세팅 완료! 테스트용으로 사운드를 한 번 재생합니다.")
	c.SE.Play()

	// Setup Crawlers
	c.Crawlers = make([]crawlers.Crawler, 0)
	if c.NewsCfg.Asiae.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Asiae{})
		if !c.NewsCfg.Asiae.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Asiae{}.GetName())
		}
	}
	if c.NewsCfg.Edaily.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Edaily{})
		if !c.NewsCfg.Edaily.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Edaily{}.GetName())
		}
	}
	if c.NewsCfg.Etoday.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Etoday{})
		if !c.NewsCfg.Etoday.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Etoday{}.GetName())
		}
	}
	if c.NewsCfg.MT.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MT{})
		if !c.NewsCfg.MT.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MT{}.GetName())
		}
	}
	if c.NewsCfg.Sedaily.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Sedaily{})
		if !c.NewsCfg.Sedaily.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Sedaily{}.GetName())
		}
	}
	if c.NewsCfg.BizChosun.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.BizChosun{})
		if !c.NewsCfg.BizChosun.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.BizChosun{}.GetName())
		}
	}
	if c.NewsCfg.FnNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.FnNews{})
		if !c.NewsCfg.FnNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.FnNews{}.GetName())
		}
	}
	if c.NewsCfg.Hankyung.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Hankyung{})
		if !c.NewsCfg.Hankyung.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Hankyung{}.GetName())
		}
	}
	if c.NewsCfg.InfoStockDaily.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.InfoStockDaily{})
		if !c.NewsCfg.InfoStockDaily.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.InfoStockDaily{}.GetName())
		}
	}
	if c.NewsCfg.MK.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MK{})
		if !c.NewsCfg.MK.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MK{}.GetName())
		}
	}
	if c.NewsCfg.MTN.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MTN{})
		if !c.NewsCfg.MTN.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MTN{}.GetName())
		}
	}
	if c.NewsCfg.Newspim.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Newspim{})
		if !c.NewsCfg.Newspim.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Newspim{}.GetName())
		}
	}
	if c.NewsCfg.YNA.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.YNA{})
		if !c.NewsCfg.YNA.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.YNA{}.GetName())
		}
	}
	/* Group 3 */
	if c.NewsCfg.BioSpectator.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.BioSpectator{})
		if !c.NewsCfg.BioSpectator.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.BioSpectator{}.GetName())
		}
	}
	if c.NewsCfg.DailyMedi.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DailyMedi{})
		if !c.NewsCfg.DailyMedi.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DailyMedi{}.GetName())
		}
	}
	if c.NewsCfg.DocDocDoc.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Doc3{})
		if !c.NewsCfg.DocDocDoc.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Doc3{}.GetName())
		}
	}
	if c.NewsCfg.DoctorsNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DoctorsNews{})
		if !c.NewsCfg.DoctorsNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DoctorsNews{}.GetName())
		}
	}
	if c.NewsCfg.MDToday.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MDToday{})
		if !c.NewsCfg.MDToday.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MDToday{}.GetName())
		}
	}
	if c.NewsCfg.News1.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.News1{})
		if !c.NewsCfg.News1.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.News1{}.GetName())
		}
	}
	if c.NewsCfg.Newsis.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Newsis{})
		if !c.NewsCfg.Newsis.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Newsis{}.GetName())
		}
	}
	if c.NewsCfg.NewsRun.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsRun{})
		if !c.NewsCfg.NewsRun.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsRun{}.GetName())
		}
	}
	if c.NewsCfg.PaxnetNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.PaxnetNews{})
		if !c.NewsCfg.PaxnetNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.PaxnetNews{}.GetName())
		}
	}
	if c.NewsCfg.Yakup.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Yakup{})
		if !c.NewsCfg.Yakup.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Yakup{}.GetName())
		}
	}
	/* Group 4 */
	if c.NewsCfg.BusinessPost.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.BusinessPost{})
		if !c.NewsCfg.BusinessPost.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.BusinessPost{}.GetName())
		}
	}
	if c.NewsCfg.DDaily.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DDaily{})
		if !c.NewsCfg.DDaily.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DDaily{}.GetName())
		}
	}
	if c.NewsCfg.DT.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DT{})
		if !c.NewsCfg.DT.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DT{}.GetName())
		}
	}
	if c.NewsCfg.GENews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.GENews{})
		if !c.NewsCfg.GENews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.GENews{}.GetName())
		}
	}
	if c.NewsCfg.INews24.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.INews24{})
		if !c.NewsCfg.INews24.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.INews24{}.GetName())
		}
	}
	if c.NewsCfg.InTheNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.InTheNews{})
		if !c.NewsCfg.InTheNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.InTheNews{}.GetName())
		}
	}
	if c.NewsCfg.Medipana.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Medipana{})
		if !c.NewsCfg.Medipana.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Medipana{}.GetName())
		}
	}
	if c.NewsCfg.Newsway.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Newsway{})
		if !c.NewsCfg.Newsway.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Newsway{}.GetName())
		}
	}
	if c.NewsCfg.Nspna.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Nspna{})
		if !c.NewsCfg.Nspna.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Nspna{}.GetName())
		}
	}
	if c.NewsCfg.SeoulWire.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.SeoulWire{})
		if !c.NewsCfg.SeoulWire.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.SeoulWire{}.GetName())
		}
	}
	if c.NewsCfg.TheBell.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.TheBell{})
		if !c.NewsCfg.TheBell.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.TheBell{}.GetName())
		}
	}
	/* Group 5 */
	if c.NewsCfg.NewsPrime.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsPrime{})
		if !c.NewsCfg.NewsPrime.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsPrime{}.GetName())
		}
	}
	if c.NewsCfg.PaxeTV.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.PaxeTV{})
		if !c.NewsCfg.PaxeTV.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.PaxeTV{}.GetName())
		}
	}
	if c.NewsCfg.DailyPharm.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DailyPharm{})
		if !c.NewsCfg.DailyPharm.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DailyPharm{}.GetName())
		}
	}
	if c.NewsCfg.SedailyGA05.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.SedailyGA05{})
		if !c.NewsCfg.SedailyGA05.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.SedailyGA05{}.GetName())
		}
	}
	if c.NewsCfg.SedailyGA07.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.SedailyGA07{})
		if !c.NewsCfg.SedailyGA07.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.SedailyGA07{}.GetName())
		}
	}
	if c.NewsCfg.RPM9.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.RPM9{})
		if !c.NewsCfg.RPM9.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.RPM9{}.GetName())
		}
	}
	if c.NewsCfg.MediaPen.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MediaPen{})
		if !c.NewsCfg.MediaPen.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MediaPen{}.GetName())
		}
	}
	if c.NewsCfg.GameFocus.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.GameFocus{})
		if !c.NewsCfg.GameFocus.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.GameFocus{}.GetName())
		}
	}
	if c.NewsCfg.MTMoneys.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MTMoneyS{})
		if !c.NewsCfg.MTMoneys.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MTMoneyS{}.GetName())
		}
	}
	/* Group 6 */
	if c.NewsCfg.Nspna11.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Nspna11{})
		if !c.NewsCfg.Nspna11.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Nspna11{}.GetName())
		}
	}
	if c.NewsCfg.Nspna21.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Nspna21{})
		if !c.NewsCfg.Nspna21.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Nspna21{}.GetName())
		}
	}
	if c.NewsCfg.NewsPrime57.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsPrime57{})
		if !c.NewsCfg.NewsPrime57.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsPrime57{}.GetName())
		}
	}
	/* Group 7 */
	if c.NewsCfg.NewsPrime67.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsPrime67{})
		if !c.NewsCfg.NewsPrime67.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsPrime67{}.GetName())
		}
	}
	if c.NewsCfg.CEOScoreDaily.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.CEOScoreDaily{})
		if !c.NewsCfg.CEOScoreDaily.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.CEOScoreDaily{}.GetName())
		}
	}
	if c.NewsCfg.ETNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.ETNews{})
		if !c.NewsCfg.ETNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.ETNews{}.GetName())
		}
	}
	if c.NewsCfg.KmedInfo.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.KmedInfo{})
		if !c.NewsCfg.KmedInfo.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.KmedInfo{}.GetName())
		}
	}
	if c.NewsCfg.Viva100.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Viva100{})
		if !c.NewsCfg.Viva100.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Viva100{}.GetName())
		}
	}
	if c.NewsCfg.ZDNet.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.ZDNet{})
		if !c.NewsCfg.ZDNet.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.ZDNet{}.GetName())
		}
	}
	/* Group 8 */
	if c.NewsCfg.AjuNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.AjuNews{})
		if !c.NewsCfg.AjuNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.AjuNews{}.GetName())
		}
	}
	if c.NewsCfg.EBN.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.EBN{})
		if !c.NewsCfg.EBN.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.EBN{}.GetName())
		}
	}
	if c.NewsCfg.KMIB.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.KMIB{})
		if !c.NewsCfg.KMIB.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.KMIB{}.GetName())
		}
	}
	if c.NewsCfg.MedicalTimes.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MedicalTimes{})
		if !c.NewsCfg.MedicalTimes.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MedicalTimes{}.GetName())
		}
	}
	if c.NewsCfg.TF.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.TF{})
		if !c.NewsCfg.TF.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.TF{}.GetName())
		}
	}
	if c.NewsCfg.GameFocus22r09.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.GameFocus22r09{})
		if !c.NewsCfg.GameFocus22r09.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.GameFocus22r09{}.GetName())
		}
	}
	/* Group 9 */
	if c.NewsCfg.LawIssue.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.LawIssue{})
		if !c.NewsCfg.LawIssue.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.LawIssue{}.GetName())
		}
	}
	if c.NewsCfg.YouthDaily.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.YouthDaily{})
		if !c.NewsCfg.YouthDaily.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.YouthDaily{}.GetName())
		}
	}
	if c.NewsCfg.KukiNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.KukiNews{})
		if !c.NewsCfg.KukiNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.KukiNews{}.GetName())
		}
	}
	if c.NewsCfg.WowTV.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.WowTV{})
		if !c.NewsCfg.WowTV.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.WowTV{}.GetName())
		}
	}
	if c.NewsCfg.NewsPrimeYMH.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsPrimeYMH{})
		if !c.NewsCfg.NewsPrimeYMH.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsPrimeYMH{}.GetName())
		}
	}
	/* Group 10 */
	if c.NewsCfg.GetNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.GetNews{})
		if !c.NewsCfg.GetNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.GetNews{}.GetName())
		}
	}
	if c.NewsCfg.NewsTown.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsTown{})
		if !c.NewsCfg.NewsTown.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsTown{}.GetName())
		}
	}
	if c.NewsCfg.DealSite.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DealSite{})
		if !c.NewsCfg.DealSite.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DealSite{}.GetName())
		}
	}
	if c.NewsCfg.PharmStock.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.PharmStock{})
		if !c.NewsCfg.PharmStock.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.PharmStock{}.GetName())
		}
	}
	if c.NewsCfg.Press9.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Press9{})
		if !c.NewsCfg.Press9.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Press9{}.GetName())
		}
	}
	/* Group 11 */
	if c.NewsCfg.Kiwoom.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Kiwoom{})
		if !c.NewsCfg.Kiwoom.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Kiwoom{}.GetName())
		}
	}
	if c.NewsCfg.HankyungBio.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.HankyungBio{})
		if !c.NewsCfg.HankyungBio.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.HankyungBio{}.GetName())
		}
	}
	if c.NewsCfg.HankyungMarketInsight.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.HankyungMarketInsight{})
		if !c.NewsCfg.HankyungMarketInsight.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.HankyungMarketInsight{}.GetName())
		}
	}
	if c.NewsCfg.BeyondPost.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.BeyondPost{})
		if !c.NewsCfg.BeyondPost.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.BeyondPost{}.GetName())
		}
	}
	if c.NewsCfg.TheGuru.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.TheGuru{})
		if !c.NewsCfg.TheGuru.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.TheGuru{}.GetName())
		}
	}
	if c.NewsCfg.NewsWorks.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.NewsWorks{})
		if !c.NewsCfg.NewsWorks.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.NewsWorks{}.GetName())
		}
	}
	if c.NewsCfg.Econovill.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Econovill{})
		if !c.NewsCfg.Econovill.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Econovill{}.GetName())
		}
	}
	if c.NewsCfg.DNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DNews{})
		if !c.NewsCfg.DNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DNews{}.GetName())
		}
	}
	if c.NewsCfg.CCReview.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.CCReview{})
		if !c.NewsCfg.CCReview.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.CCReview{}.GetName())
		}
	}
	if c.NewsCfg.TheElec.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.TheElec{})
		if !c.NewsCfg.TheElec.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.TheElec{}.GetName())
		}
	}
	if c.NewsCfg.News1Latest.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.News1Latest{})
		if !c.NewsCfg.News1Latest.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.News1Latest{}.GetName())
		}
	}
	if c.NewsCfg.MKVIP26.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MKVIP26{})
		if !c.NewsCfg.MKVIP26.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MKVIP26{}.GetName())
		}
	}
	if c.NewsCfg.MKVIP10001.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MKVIP10001{})
		if !c.NewsCfg.MKVIP10001.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MKVIP10001{}.GetName())
		}
	}
	if c.NewsCfg.Etoday1202.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Etoday1202{})
		if !c.NewsCfg.Etoday1202.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Etoday1202{}.GetName())
		}
	}
	/* Group 12 */
	if c.NewsCfg.MyAsset.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MyAsset{})
		if !c.NewsCfg.MyAsset.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MyAsset{}.GetName())
		}
	}
	if c.NewsCfg.HeraldCorp.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.HeraldCorp{})
		if !c.NewsCfg.HeraldCorp.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.HeraldCorp{}.GetName())
		}
	}
	if c.NewsCfg.Bosa.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Bosa{})
		if !c.NewsCfg.Bosa.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Bosa{}.GetName())
		}
	}
	if c.NewsCfg.HITNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.HITNews{})
		if !c.NewsCfg.HITNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.HITNews{}.GetName())
		}
	}
	if c.NewsCfg.DataNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DataNews{})
		if !c.NewsCfg.DataNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DataNews{}.GetName())
		}
	}
	if c.NewsCfg.DoctorsTimes.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.DoctorsTimes{})
		if !c.NewsCfg.DoctorsTimes.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.DoctorsTimes{}.GetName())
		}
	}
	if c.NewsCfg.Whosaeng.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.Whosaeng{})
		if !c.NewsCfg.Whosaeng.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.Whosaeng{}.GetName())
		}
	}
	if c.NewsCfg.HealthInNews.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.HealthInNews{})
		if !c.NewsCfg.HealthInNews.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.HealthInNews{}.GetName())
		}
	}
	if c.NewsCfg.MedipharmHealth.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.MedipharmHealth{})
		if !c.NewsCfg.MedipharmHealth.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.MedipharmHealth{}.GetName())
		}
	}
	if c.NewsCfg.AsiaeFeature.Crawl {
		c.Crawlers = append(c.Crawlers, crawlers.AsiaeFeature{})
		if !c.NewsCfg.AsiaeFeature.TitleFiltering {
			CondOnlyContentsCrawlers = append(CondOnlyContentsCrawlers, crawlers.AsiaeFeature{}.GetName())
		}
	}

	c.Logger.Println("Crawler 세팅 완료!")
	c.Logger.Printf("제목 필터링 기능이 꺼진 크롤러들 : [%s]\n", strings.Join(CondOnlyContentsCrawlers, ","))

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
				for idx := range data {
					if !utils.IsContain(data[idx].URL, prevData[name]) { // New Item
						// Exceptions (do not check title/keyword condition)
						if utils.IsContain(crawler.GetName(), CondExceptionCrawlers) {
							err = crawler.GetContents(&data[idx])
							if err != nil {
								c.Logger.Printf("[ERROR] crawler.GetContents() : crawler(%s) : idx(%d) : %s", name, idx, err)
								errorCount[name]++
								continue
							}
							// Messaging Goroutine
							go c.TG.SendMessage(data[idx], []string{""})
							go c.SE.Play()
							continue
						}

						// if TitleCond true (or if the crawler is in CondOnlyContentsCrawlers)
						if utils.TitleCond(data[idx]) || utils.IsContain(crawler.GetName(), CondOnlyContentsCrawlers) {
							err = crawler.GetContents(&data[idx])
							if err != nil {
								c.Logger.Printf("[ERROR] crawler.GetContents() : crawler(%s) : idx(%d) : %s", name, idx, err)
								errorCount[name]++
								continue
							}

							detectKeywords, ok := utils.KeywordCond(data[idx], c.Cfg.Keywords, c.Cfg.Filters)
							if ok && len(detectKeywords) >= c.Cfg.Crawler.KeywordDetectionNum {
								go c.TG.SendMessage(data[idx], detectKeywords)
								go c.SE.Play()
								continue
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
