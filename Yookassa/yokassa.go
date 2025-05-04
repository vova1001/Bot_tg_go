package yookassa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	// Маршалинг данных в JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге JSON: %v", err)
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/webhooks", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(shopID, secretKey)

	// Отправка запроса
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении тела ответа: %v", err)
	}

	// Проверка статуса ответа
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("не удалось зарегистрировать вебхук, статус: %d, ответ: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
