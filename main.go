package main

import (
	"CourseTg/config"
	"CourseTg/internal/handlers"
	"CourseTg/webhook"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg := config.LoadConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
	}
	bot.Debug = true
	log.Printf("✅ Авторизован как %s", bot.Self.UserName)

	webhookURL := "https://bot-tg-go.onrender.com"
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatal("Ошибка при установке вебхука:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var update tgbotapi.Update
			if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			handlers.HandleUpdates(bot, update)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/webhook/yookassa", webhook.HandleYooKassaWebhook)

	go func() {
		for {
			resp, err := http.Get("https://bot-tg-go.onrender.com/ping")
			if err != nil {
				log.Println("🔁 Ошибка пинга:", err)
			} else {
				log.Println("✅ Сервер активен (ping)")
				resp.Body.Close()
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
