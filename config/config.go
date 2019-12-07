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
	DelayTimer          int64  `json:"delay_timer"`
	MaxProcs            int    `json:"max_procs"`
	KeywordDetectionNum int    `json:"keyword_detection_num"`
}

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	Crawler  CrawlerConfig  `json:"crawler"`
	Keywords []string
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
