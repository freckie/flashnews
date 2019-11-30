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
	Asiae   bool `json:"asiae"`
	Edaily  bool `json:"edaily"`
	Etoday  bool `json:"etoday"`
	MT      bool `json:"mt"`
	Sedaily bool `json:"sedaily"`
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
