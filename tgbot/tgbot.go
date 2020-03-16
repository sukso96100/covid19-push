package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sukso96100/covid19-push/database"
)

type TgChanBot struct {
	Bot     *tgbotapi.BotAPI
	Channel string
}

var botObj *TgChanBot

func InitTgBot(token string, channel string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	botObj = &TgChanBot{
		Bot:     bot,
		Channel: channel,
	}
	return nil
}

func Bot() *TgChanBot {
	return botObj
}

func (bot *TgChanBot) SendStatMsg(current database.StatData) error {
	msgContent := fmt.Sprintf(
		"코로나 19 발생 현황\n%s\n",
		database.CreateStatMsg(current))

	msg := tgbotapi.NewMessageToChannel(bot.Channel, msgContent)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("발생현황 자세히 보기", "http://ncov.mohw.go.kr/index.jsp"),
		),
	)
	if _, err := botObj.Bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (bot *TgChanBot) SendNewsMsg(newsData database.NewsData) error {
	msgContent := fmt.Sprintf(
		"%s\n(담당부서: %s)",
		newsData.Title, newsData.Department)
	msg := tgbotapi.NewMessageToChannel(bot.Channel, msgContent)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("자료 자세히 보기", newsData.Link),
		),
	)
	if _, err := botObj.Bot.Send(msg); err != nil {
		return err // Again, this is a bad way to handle errors.
	}
	return nil
}
