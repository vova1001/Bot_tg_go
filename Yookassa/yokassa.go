package yookassa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func RegisterWebhook() error {
	shopID := os.Getenv("YOOKASSA_SHOP_ID")
	secretKey := os.Getenv("YOOKASSA_SECRET_KEY")

	// Создание payload для регистрации вебхука
	payload := map[string]interface{}{
		"event": "payment.succeeded",
		"url":   "https://bot-tg-go.onrender.com/yookassa-webhook",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/webhooks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(shopID, secretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("не удалось зарегистрировать вебхук, статус: %d", resp.StatusCode)
	}

	return nil
}
