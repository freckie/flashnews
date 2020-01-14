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

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	Crawler  CrawlerConfig  `json:"crawler"`
	Keywords []string
	Filters  []string
}

type NewsConfig struct {
	Asiae          bool `json:"asiae.co.kr"`
	Edaily         bool `json:"edaily.co.kr"`
	Etoday         bool `json:"etoday.co.kr"`
	MT             bool `json:"mt.co.kr"`
	Sedaily        bool `json:"sedaily.com"`
	BizChosun      bool `json:"biz.chosun.com"`
	FnNews         bool `json:"fnnews.com"`
	Hankyung       bool `json:"hankyung.com"`
	InfoStockDaily bool `json:"infostockdaily.co.kr"`
	MK             bool `json:"mk.co.kr"`
	MTN            bool `json:"mtn.co.kr"`
	Newspim        bool `json:"newspim.com"`
	YNA            bool `json:"yna.co.kr"`
	/* Group 3 */
	BioSpectator bool `json:"biospectator.com"`
	DailyMedi    bool `json:"dailymedi.com"`
	DocDocDoc    bool `json:"docdocdoc.co.kr"`
	DoctorsNews  bool `json:"doctorsnews.co.kr"`
	MDToday      bool `json:"mdtoday.co.kr"`
	News1        bool `json:"news1.kr"`
	Newsis       bool `json:"newsis.com"`
	NewsRun      bool `json:"newsrun.co.kr"`
	PaxnetNews   bool `json:"paxnetnews.com"`
	Yakup        bool `json:"yakup.com"`
	/* Group 4 */
	BusinessPost bool `json:"businesspost.co.kr"`
	DDaily       bool `json:"ddaily.co.kr"`
	DT           bool `json:"dt.co.kr"`
	GENews       bool `json:"g-enews.com"`
	INews24      bool `json:"inews24.com"`
	InTheNews    bool `json:"inthenews.co.kr"`
	Medipana     bool `json:"medipana.com"`
	Newsway      bool `json:"newsway.co.kr"`
	Nspna        bool `json:"nspna.com"`
	SeoulWire    bool `json:"seoulwire.com"`
	TheBell      bool `json:"thebell.co.kr"`
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
