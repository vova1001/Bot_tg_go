package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MyBot — обёртка над tgbotapi.BotAPI (по желанию, можно не оборачивать)
type MyBot struct {
	*tgbotapi.BotAPI
}

// NewBotAPI инициализирует нового бота
func NewBotAPI(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Бот авторизован: %s", bot.Self.UserName)
	return bot, nil
}
