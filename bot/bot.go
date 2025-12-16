package bot

import (
	"log"
	"s21/service"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Run() {
	service := service.NewMySerivce()
	bot, updates := botInitMust()
	maxLength := 4000

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// update.Message.Text

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			var text string
			switch update.Message.Command() {
			case "start":
				text, _ = All4MskTribesLogins(service)
			case "end":
				text = service.Usernames()
			case "seat":
				text, _ = PartisipantSeatInfo(*service, strings.Fields(string(update.Message.Text))[1])
			case "survivors":
				text, _ = Survivors(service)
			}

			for len(text) > 0 {
				chunk := text
				if len(chunk) > maxLength {
					chunk = text[:maxLength]
					text = text[maxLength:]
				} else {
					text = ""
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, chunk)
				bot.Send(msg)
			}
		}
	}
}

func botInitMust() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(tokenBot)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Авторизован как %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Println(err)
	}

	return bot, updates
}
