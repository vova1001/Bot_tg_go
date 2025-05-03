package payment

import (
	"bytes"
	"encoding/json"
	"net/http"
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

func CreatePayment(amount, description, telegramID, shopID, secretKey string) (string, error) {
	url := "https://api.yookassa.ru/v3/payments"

	reqBody := PaymentRequest{}
	reqBody.Amount.Value = amount
	reqBody.Amount.Currency = "RUB"
	reqBody.Confirmation.Type = "redirect"
	reqBody.Confirmation.ReturnURL = "https://t.me/JuliiaFitness_bot" // Измени!
	reqBody.Capture = true
	reqBody.Description = description
	reqBody.Metadata.TelegramID = telegramID

	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(shopID, secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res PaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.Confirmation.ConfirmationURL, nil
}
