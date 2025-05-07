package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type Item struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Amount      Amount  `json:"amount"`
	VATCode     int     `json:"vat_code"`
}

type Customer struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type Receipt struct {
	Customer Customer `json:"customer"`
	Items    []Item   `json:"items"`
}

type Metadata struct {
	TelegramID string `json:"telegram_id"`
	CourseID   string `json:"course_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type PaymentRequest struct {
	Amount       Amount       `json:"amount"`
	Confirmation Confirmation `json:"confirmation"`
	Capture      bool         `json:"capture"`
	Description  string       `json:"description"`
	Metadata     Metadata     `json:"metadata"`
	Receipt      Receipt      `json:"receipt"`
}

type PaymentResponse struct {
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func CreatePayment(amount, description, telegramID, courseID, name, email, shopID, secretKey string) (string, error) {
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
			Name:       name,
			Email:      email,
		},
		Receipt: Receipt{
			Customer: Customer{
				Email:    email,
				FullName: name,
			},
			Items: []Item{
				{
					Description: description,
					Quantity:    1,
					Amount: Amount{
						Value:    amount,
						Currency: "RUB",
					},
					VATCode: 1,
				},
			},
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
	req.Header.Set("Idempotence-Key", uuid.New().String())

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
