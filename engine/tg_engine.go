package engine

import (
	"errors"
	"log"
	"strings"

	"flashnews/config"
	"flashnews/models"
	"flashnews/utils"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const MaxPrevMessageQueueSize = 5

type TGEngine struct {
	Bot          *telegram.BotAPI
	Cfg          *config.Config
	PrevMessages []string
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
	if tg.IsDuplicated(item) {
		return errors.New("Message Duplicated.")
	}
	tg.AddMessage(item)

	keywordStr := "[" + strings.Join(keywords, ", ") + "]"
	contentsStr := strings.Replace(utils.StringSplit(item.Contents, 300), "<", "", -1)
	contentsStr = strings.Replace(contentsStr, ">", "", -1)

	msgStr := tg.Cfg.Telegram.MessageFormat
	msgStr = strings.Replace(msgStr, "%(title)", "<b>"+item.Title+"</b>", -1)
	msgStr = strings.Replace(msgStr, "%(contents)", contentsStr, -1)
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

func (tg TGEngine) IsDuplicated(item models.NewsItem) bool {
	for _, prevMsg := range tg.PrevMessages {
		if item.Title == prevMsg {
			return true
		}
	}
	return false
}

func (tg *TGEngine) AddMessage(item models.NewsItem) {
	temp := make([]string, MaxPrevMessageQueueSize)

	for i := 0; i < MaxPrevMessageQueueSize; i++ {
		temp[i+1] = tg.PrevMessages[i]
	}
	temp[0] = item.Title

	tg.PrevMessages = temp
}
