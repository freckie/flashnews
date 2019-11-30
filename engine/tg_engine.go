package engine

import (
	"log"
	"strings"

	"flashnews/config"
	"flashnews/models"
	"flashnews/utils"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TGEngine struct {
	Bot *telegram.BotAPI
	Cfg *config.Config
}

func (tg *TGEngine) GenerateBot() error {
	bot, err := telegram.NewBotAPI(tg.Cfg.Telegram.BotToken)
	if err != nil {
		return err
	}

	tg.Bot = bot
	return nil
}

func (tg TGEngine) SendMessage(item models.NewsItem, keywords []string) error {
	keywordStr := "[" + strings.Join(keywords, ", ") + "]"

	msgStr := tg.Cfg.Telegram.MessageFormat
	msgStr = strings.Replace(msgStr, "%(title)", "<b>"+item.Title+"</b>", -1)
	msgStr = strings.Replace(msgStr, "%(contents)", utils.StringSplit(item.Contents, 300), -1)
	msgStr = strings.Replace(msgStr, "%(keywords)", keywordStr, -1)
	msgStr = strings.Replace(msgStr, "%(link)", item.URL, -1)

	for _, channel := range tg.Cfg.Telegram.Channels {
		msgType := telegram.MessageConfig{
			BaseChat: telegram.BaseChat{
				ChatID:           channel,
				ReplyToMessageID: 0,
			},
			Text:                  msgStr,
			ParseMode:             "html",
			DisableWebPagePreview: false,
		}
		sentMsg, err := tg.Bot.Send(msgType)
		if err != nil {
			log.Println("[ERROR] 메세지 전송 실패 : ", err)
		}
		log.Printf("채널(%d)에 메세지 전송 : %v", channel, sentMsg.Text)
	}

	return nil
}

func (tg TGEngine) TestMessage() error {
	msgStr := "<b>테스트 메세지입니다.</b>"

	for _, channel := range tg.Cfg.Telegram.Channels {
		msgType := telegram.MessageConfig{
			BaseChat: telegram.BaseChat{
				ChatID:           channel,
				ReplyToMessageID: 0,
			},
			Text:                  msgStr,
			ParseMode:             "html",
			DisableWebPagePreview: false,
		}
		sentMsg, err := tg.Bot.Send(msgType)
		if err != nil {
			log.Println("[ERROR] 테스트 메세지 전송 실패 : ", err)
		}
		log.Printf("채널(%d)에 테스트 메세지 전송 : %v", channel, sentMsg.Text)
	}

	return nil
}
