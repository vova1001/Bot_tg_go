package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MyBot struct {
	*tgbotapi.BotAPI
}

func NewBotAPI(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Бот авторизован: %s", bot.Self.UserName)
	return bot, nil
}
