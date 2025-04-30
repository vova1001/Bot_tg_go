package webhook

import (
	"encoding/json"
	"log"
	"net/http"
)

type YooKassaNotification struct {
	// Структура для уведомлений от YooKassa
	Event  string `json:"event"`
	Object struct {
		Status string `json:"status"`
	} `json:"object"`
}

func HandleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var notification YooKassaNotification
		err := json.NewDecoder(r.Body).Decode(&notification)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Пример обработки: логирование статуса платежа
		log.Printf("Received notification: event=%s, status=%s", notification.Event, notification.Object.Status)

		// Ответ в YooKassa, подтверждающий получение
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
