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
	InputPath  string `json:"input_path"`
	DelayTimer int64  `json:"delay_timer"`
	DelayAuto  bool
}

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	Crawler  CrawlerConfig  `json:"crawler"`
	Keywords []string
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := &Config{}

	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	json.Unmarshal(dataBytes, cfg)

	if cfg.Crawler.DelayTimer < 0 {
		cfg.Crawler.DelayAuto = true
	} else {
		cfg.Crawler.DelayAuto = false
	}

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
