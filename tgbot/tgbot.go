package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sukso96100/covid19-push/database"
)

type TgChanBot struct{
	Bot *tgbotapi.BotAPI
	Channel string
} 
var botObj *TgChanBot 

func InitTgBot(token string, channel string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	botObj = &TgChanBot{
		Bot: bot,
		Channel: channel,
	}
	return nil
}

func Bot() *TgChanBot {
	return botObj
}

func (bot *TgChanBot) SendStatMsg(prev database.StatData, current database.StatData) error {
	msgContent := fmt.Sprintf(
		"코로나 19 발생 현황\n%s\nhttp://ncov.mohw.go.kr/bdBoardList_Real.do",
		database.CreateStatMsg(prev, current))

	msg := tgbotapi.NewMessageToChannel(bot.Channel, msgContent)
	if _, err := botObj.Bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (bot *TgChanBot) SendNewsMsg(newsData database.NewsData) error {
	msgContent := fmt.Sprintf(
		"%s\n(담당부서: %s)\n%s",
		newsData.Title, newsData.Department, newsData.Link)
	msg := tgbotapi.NewMessageToChannel(bot.Channel, msgContent)
	if _, err := botObj.Bot.Send(msg); err != nil {
		return err // Again, this is a bad way to handle errors.
	}
	return nil
}