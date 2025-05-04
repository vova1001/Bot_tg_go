package config

import (
	"log"
	"os"
)

type Config struct {
	TelegramToken     string
	YooKassaShopID    string
	YooKassaSecretKey string
}

func LoadConfig() *Config {
	return &Config{
		TelegramToken:     getEnv("TELEGRAM_TOKEN", ""),
		YooKassaShopID:    getEnv("YOO_KASSA_SHOP_ID", ""),
		YooKassaSecretKey: getEnv("YOO_KASSA_SECRET_KEY", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == "" {
			log.Fatalf("❌ Ошибка: переменная окружения %s не установлена", key)
		}
		return defaultValue
	}
	return value
}
