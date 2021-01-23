package config

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type TelegramConfig struct {
	BotToken      string  `json:"bot_token"`
	Channels      []int64 `json:"channels"`
	MessageFormat string  `json:"message_format"`
}

type CrawlerConfig struct {
	InputPath           string `json:"input_path"`
	InputPath2          string `json:"input_path2"`
	InputPath3          string `json:"input_path3"`
	DelayTimer          int64  `json:"delay_timer"`
	MaxProcs            int    `json:"max_procs"`
	KeywordDetectionNum int    `json:"keyword_detection_num"`
}

type SoundConfig struct {
	On       bool   `json:"on"`
	FilePath string `json:"file_path"`
}

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	Crawler  CrawlerConfig  `json:"crawler"`
	Sound    SoundConfig    `json:"sound"`
	Keywords []string
	Filters  []string
}

type NewsConfig struct {
	Asiae          NewsConfigItem `json:"asiae.co.kr"`
	Edaily         NewsConfigItem `json:"edaily.co.kr"`
	Etoday         NewsConfigItem `json:"etoday.co.kr"`
	MT             NewsConfigItem `json:"mt.co.kr"`
	Sedaily        NewsConfigItem `json:"sedaily.com"`
	BizChosun      NewsConfigItem `json:"biz.chosun.com"`
	FnNews         NewsConfigItem `json:"fnnews.com"`
	Hankyung       NewsConfigItem `json:"hankyung.com"`
	InfoStockDaily NewsConfigItem `json:"infostockdaily.co.kr"`
	MK             NewsConfigItem `json:"mk.co.kr"`
	MTN            NewsConfigItem `json:"mtn.co.kr"`
	Newspim        NewsConfigItem `json:"newspim.com"`
	YNA            NewsConfigItem `json:"yna.co.kr"`
	/* Group 3 */
	BioSpectator NewsConfigItem `json:"biospectator.com"`
	DailyMedi    NewsConfigItem `json:"dailymedi.com"`
	DocDocDoc    NewsConfigItem `json:"docdocdoc.co.kr"`
	DoctorsNews  NewsConfigItem `json:"doctorsnews.co.kr"`
	MDToday      NewsConfigItem `json:"mdtoday.co.kr"`
	News1        NewsConfigItem `json:"news1.kr"`
	Newsis       NewsConfigItem `json:"newsis.com"`
	NewsRun      NewsConfigItem `json:"newsrun.co.kr"`
	PaxnetNews   NewsConfigItem `json:"paxnetnews.com"`
	Yakup        NewsConfigItem `json:"yakup.com"`
	/* Group 4 */
	BusinessPost NewsConfigItem `json:"businesspost.co.kr"`
	DDaily       NewsConfigItem `json:"ddaily.co.kr"`
	DT           NewsConfigItem `json:"dt.co.kr"`
	GENews       NewsConfigItem `json:"g-enews.com"`
	INews24      NewsConfigItem `json:"inews24.com"`
	InTheNews    NewsConfigItem `json:"inthenews.co.kr"`
	Medipana     NewsConfigItem `json:"medipana.com"`
	Newsway      NewsConfigItem `json:"newsway.co.kr"`
	Nspna        NewsConfigItem `json:"nspna.com"`
	SeoulWire    NewsConfigItem `json:"seoulwire.com"`
	TheBell      NewsConfigItem `json:"thebell.co.kr"`
	/* Group 5 */
	NewsPrime   NewsConfigItem `json:"newsprime.co.kr"`
	PaxeTV      NewsConfigItem `json:"paxetv.com"`
	DailyPharm  NewsConfigItem `json:"dailypharm.com"`
	SedailyGA05 NewsConfigItem `json:"sedaily.com/NewsList/GA05"`
	SedailyGA07 NewsConfigItem `json:"sedaily.com/NewsList/GA07"`
	RPM9        NewsConfigItem `json:"rpm9.com"`
	MediaPen    NewsConfigItem `json:"mediapen.com"`
	GameFocus   NewsConfigItem `json:"gamefocus.co.kr"`
	MTMoneys    NewsConfigItem `json:"moneys.mt.co.kr"`
	/* Group 6 */
	Nspna11     NewsConfigItem `json:"nspna.com/news/?cid=11"`
	Nspna21     NewsConfigItem `json:"nspna.com/news/?cid=21"`
	NewsPrime57 NewsConfigItem `json:"newsprime.co.kr/section_list_all/?sec_no=57"`
	/* Group 7 */
	NewsPrime67   NewsConfigItem `json:"newsprime.co.kr/section_list_all/?sec_no=67"`
	CEOScoreDaily NewsConfigItem `json:"ceoscoredaily.com"`
	ETNews        NewsConfigItem `json:"etnews.com"`
	KmedInfo      NewsConfigItem `json:"kmedinfo.co.kr"`
	Viva100       NewsConfigItem `json:"viva100.com"`
	ZDNet         NewsConfigItem `json:"zdnet.co.kr"`
	/* Group 8 */
	AjuNews        NewsConfigItem `json:"ajunews.com"`
	EBN            NewsConfigItem `json:"ebn.co.kr"`
	KMIB           NewsConfigItem `json:"news.kmib.co.kr"`
	MedicalTimes   NewsConfigItem `json:"medicaltimes.com"`
	TF             NewsConfigItem `json:"news.tf.co.kr"`
	GameFocus22r09 NewsConfigItem `json:"gamefocus.co.kr/section.php?thread=22r09"`
	/* Group 9 */
	LawIssue     NewsConfigItem `json:"lawissue.co.kr"`
	YouthDaily   NewsConfigItem `json:"youthdaily.co.kr"`
	KukiNews     NewsConfigItem `json:"kukinews.com"`
	WowTV        NewsConfigItem `json:"wowtv.co.kr"`
	NewsPrimeYMH NewsConfigItem `json:"newsprime.co.kr/article_list_writer/?name=양민호+기자"`
	/* Group 10 */
	GetNews    NewsConfigItem `json:"getnews.co.kr"`
	NewsTown   NewsConfigItem `json:"newstown.co.kr"`
	DealSite   NewsConfigItem `json:"dealsite.co.kr"`
	PharmStock NewsConfigItem `json:"pharmstock.co.kr"`
	Press9     NewsConfigItem `json:"press9.kr"`
	/* Group 11 */
	Kiwoom                NewsConfigItem `json:"kiwoom.com"`
	HankyungBio           NewsConfigItem `json:"hankyung.com/bioinsight"`
	HankyungMarketInsight NewsConfigItem `json:"marketinsight.hankyung.com"`
	BeyondPost            NewsConfigItem `json:"beyondpost.co.kr"`
	TheGuru               NewsConfigItem `json:"theguru.co.kr"`
	NewsWorks             NewsConfigItem `json:"newsworks.co.kr"`
	Econovill             NewsConfigItem `json:"econovill.com"`
	DNews                 NewsConfigItem `json:"dnews.co.kr"`
	CCReview              NewsConfigItem `json:"ccreview.co.kr"`
	TheElec               NewsConfigItem `json:"thelec.kr"`
	News1Latest           NewsConfigItem `json:"news1.kr/latest"`
	MKVIP26               NewsConfigItem `json:"vip.mk.co.kr/newSt/news/news_list.php?sCode=26"`
	MKVIP10001            NewsConfigItem `json:"vip.mk.co.kr/newSt/news/news_list.php?sCode=10001"`
	Etoday1202            NewsConfigItem `json:"etoday.co.kr/news/section/subsection?MID=1202"`
	/* Group 12 */
	MyAsset         NewsConfigItem `json:"myasset.com"`
	HeraldCorp      NewsConfigItem `json:"biz.heraldcorp.com"`
	Bosa            NewsConfigItem `json:"bosa.co.kr"`
	HITNews         NewsConfigItem `json:"hitnews.co.kr"`
	DataNews        NewsConfigItem `json:"datanews.co.kr"`
	DoctorsTimes    NewsConfigItem `json:"doctorstimes.com"`
	Whosaeng        NewsConfigItem `json:"whosaeng.com"`
	HealthInNews    NewsConfigItem `json:"healthinnews.co.kr"`
	MedipharmHealth NewsConfigItem `json:"medipharmhealth.co.kr"`
}

type NewsConfigItem struct {
	Crawl          bool `json:"crawl"`
	TitleFiltering bool `json:"title_filtering"`
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := &Config{}

	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	json.Unmarshal(dataBytes, cfg)

	return cfg, nil
}

func LoadNewsConfig(filePath string) (*NewsConfig, error) {
	cfg := &NewsConfig{}

	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	json.Unmarshal(dataBytes, cfg)

	return cfg, nil
}

func LoadKeywords(filePath string) ([]string, error) {
	result := make([]string, 0)

	fo, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer fo.Close()

	reader := bufio.NewReader(fo)
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		_line := strings.Replace(string(line), " ", "", -1)
		result = append(result, _line)
	}

	return result, nil
}

func LoadFilters(filePath string) ([]string, error) {
	result := make([]string, 0)

	fo, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer fo.Close()

	reader := bufio.NewReader(fo)
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		_line := strings.Replace(string(line), " ", "", -1)
		result = append(result, _line)
	}

	return result, nil
}
