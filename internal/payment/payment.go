package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Metadata struct {
	TelegramID string `json:"telegram_id"`
}

type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Capture     bool     `json:"capture"`
	Description string   `json:"description"`
	Metadata    Metadata `json:"metadata"`
}

type PaymentResponse struct {
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func parseAmount(amountStr string) float64 {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0.0
	}
	// Округление до двух знаков после запятой
	return math.Round(amount*100) / 100
}

func CreatePayment(amount, description, telegramID, shopID, secretKey string) (string, error) {
	url := "https://api.yookassa.ru/v3/payments"

	// Создание структуры запроса
	reqBody := PaymentRequest{}
	reqBody.Amount.Value = fmt.Sprintf("%.2f", parseAmount(amount)) // Преобразуем amount в правильный формат
	reqBody.Amount.Currency = "RUB"
	reqBody.Confirmation.Type = "redirect"
	reqBody.Confirmation.ReturnURL = "https://t.me/JuliiaFitness_bot" // Убедитесь, что это правильный URL для редиректа
	reqBody.Capture = true
	reqBody.Description = description
	reqBody.Metadata.TelegramID = telegramID

	// Преобразование структуры в JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Создание HTTP запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(shopID, secretKey)

	// Уникальный Idempotence Key
	req.Header.Set("Idempotence-Key", fmt.Sprintf("key-%d", time.Now().UnixNano()))

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Обработка ответа
	var res PaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	// Возвращаем ссылку на подтверждение
	return res.Confirmation.ConfirmationURL, nil
}
