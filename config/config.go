package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken     string
	YooKassaShopID    string
	YooKassaSecretKey string
	SuccessURL        string
	FailURL           string
	WebhookURL        string // Добавлен Webhook URL
}

func LoadConfig() *Config {
	err := godotenv.Load("D:/CourseTg/.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла:", err)
	}

	return &Config{
		TelegramToken:     os.Getenv("TELEGRAM_TOKEN"),
		YooKassaShopID:    os.Getenv("YOOKASSA_SHOP_ID"),
		YooKassaSecretKey: os.Getenv("YOOKASSA_SECRET_KEY"),
		SuccessURL:        os.Getenv("SUCCESS_URL"),
		FailURL:           os.Getenv("FAIL_URL"),
		WebhookURL:        os.Getenv("WEBHOOK_URL"), // Загружаем Webhook URL
	}
}
