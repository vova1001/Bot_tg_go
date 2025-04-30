package payment

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func CreatePayment(amount string, description string, successURL string, shopID string, secretKey string) (string, error) {
	data := url.Values{}
	data.Set("amount", amount)
	data.Set("currency", "RUB")
	data.Set("description", description)
	data.Set("success_url", successURL)
	data.Set("shop_id", shopID)
	data.Set("secret_key", secretKey)

	// Формируем запрос к API Юкассы
	resp, err := http.PostForm("https://api.yookassa.ru/v3/payment", data)
	if err != nil {
		log.Println("Error while creating payment:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Проверка на успешный ответ
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create payment, status code: %d", resp.StatusCode)
	}

	// Возвращаем URL для оплаты
	return "https://payment.yookassa.ru/" + resp.Request.URL.Path, nil
}
