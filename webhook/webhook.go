package webhook

import (
	"CourseTg/config"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type YooKassaNotification struct {
	Event  string `json:"event"`
	Object struct {
		Status   string `json:"status"`
		Metadata struct {
			TelegramID string `json:"telegram_id"`
			CourseID   string `json:"course_id"` // теперь получаем и ID курса
		} `json:"metadata"`
	} `json:"object"`
}

func HandleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var notification YooKassaNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if notification.Event == "payment.succeeded" && notification.Object.Status == "succeeded" {
		log.Printf("✅ Оплата прошла для %s (%s)", notification.Object.Metadata.TelegramID, notification.Object.Metadata.CourseID)
		go sendPDF(notification.Object.Metadata.TelegramID, notification.Object.Metadata.CourseID)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func sendPDF(telegramID string, courseID string) {
	cfg := config.LoadConfig()
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Printf("❌ Ошибка создания бота: %v", err)
		return
	}

	userID, err := strconv.ParseInt(telegramID, 10, 64)
	if err != nil {
		log.Printf("❌ Неверный Telegram ID: %v", err)
		return
	}

	var filesToSend []string
	switch courseID {
	case "course_1":
		filesToSend = []string{"Kniga_receptov.pdf"}
	case "course_2":
		filesToSend = []string{"Sbornik_zavtrakov.pdf"}
	case "course_3":
		filesToSend = []string{"Kniga_receptov.pdf", "Sbornik_zavtrakov.pdf"}
	default:
		log.Printf("❌ Неизвестный courseID: %s", courseID)
		return
	}

	for _, file := range filesToSend {
		pdfPath := filepath.Join("pdf", file)
		fileBytes, err := os.ReadFile(pdfPath)
		if err != nil {
			log.Printf("❌ Ошибка чтения PDF %s: %v", file, err)
			continue
		}

		doc := tgbotapi.FileBytes{
			Name:  file,
			Bytes: fileBytes,
		}

		msg := tgbotapi.NewDocumentUpload(userID, doc)
		msg.Caption = "🎉 Спасибо за покупку! Вот ваш файл:"
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("❌ Ошибка при отправке PDF %s: %v", file, err)
		} else {
			log.Printf("📄 Файл %s отправлен пользователю %d", file, userID)
		}
	}
}
