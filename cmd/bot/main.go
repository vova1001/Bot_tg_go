package main

import (
	"CourseTg/config"
	"CourseTg/internal/handlers"
	"CourseTg/webhook" // импортируем обработчик вебхука
	"encoding/json"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg := config.LoadConfig()

	// Создаем нового бота
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	// Устанавливаем режим отладки
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Создаем вебхук URL для Telegram
	webhookURL := "https://your-domain.com/"
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем сервер для обработки запросов от Telegram
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Получаем обновление от Telegram
			var update tgbotapi.Update
			err := json.NewDecoder(r.Body).Decode(&update)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Обрабатываем обновление
			handlers.HandleUpdates(bot, update)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для уведомлений от YooKassa
	http.HandleFunc("/yookassa-webhook", webhook.HandleYooKassaWebhook)

	// Запускаем сервер
	log.Println("Server started at https://your-domain.com/")
	log.Fatal(http.ListenAndServe(":443", nil)) // Запуск HTTPS сервера
}
