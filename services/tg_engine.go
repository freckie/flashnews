package services

import (
	"log"

	"flashnews/config"
	"flashnews/models"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TGEngine struct {
	Bot *telegram.BotAPI
	Cfg *config.Config
}

func (tg TGEngine) GenerateBot() error {
	bot, err := telegram.NewBotAPI(tg.Cfg.Telegram.BotToken)
	if err != nil {
		return err
	}

	tg.Bot = bot
	return nil
}

func (tg TGEngine) SendMessage(item models.NewsItem, company string, keyword string) error {
	msgStr := tg.Cfg.Telegram.MessageFormat

	for _, channel := range tg.Cfg.Telegram.Channels {
		newMsg := telegram.NewMessage(channel, msgStr)
		sentMsg, err := tg.Bot.Send(newMsg)
		if err != nil {
			log.Println("[ERROR] 메세지 전송 실패 : ", err)
		}
		log.Printf("채널(%d)에 메세지 전송 : %v", channel, sentMsg.Text)
	}

	return nil
}
