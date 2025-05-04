package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type PaymentRequest struct {
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
	Capture      bool         `json:"capture"`
	Description  string       `json:"description"`
	Metadata     Metadata     `json:"metadata"`
}

type Metadata struct {
	TelegramID string `json:"telegram_id"`
	CourseID   string `json:"course_id"` // Новое поле
}

type PaymentResponse struct {
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func CreatePayment(amount, description, telegramID, courseID, shopID, secretKey string) (string, error) {
	requestBody := PaymentRequest{
		Amount: Amount{
			Value:    amount,
			Currency: "RUB",
		},
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: "https://t.me/JuliiaFitness_bot",
		},
		Capture:     true,
		Description: description,
		Metadata: Metadata{
			TelegramID: telegramID,
			CourseID:   courseID,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка при маршалинге запроса: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(shopID, secretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	var paymentResp PaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&paymentResp)
	if err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return paymentResp.Confirmation.ConfirmationURL, nil
}
