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
		log.Printf("✅ Оплата прошла для %s", notification.Object.Metadata.TelegramID)
		go sendPDF(notification.Object.Metadata.TelegramID)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func sendPDF(telegramID string) {
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

	// Указываем путь к PDF (на Render PDF должен быть рядом с бинарником или в отдельной папке)
	pdfPath := filepath.Join("files", "guide.pdf") // Подразумевается папка /files и файл guide.pdf

	fileBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Printf("❌ Ошибка чтения PDF: %v", err)
		return
	}

	doc := tgbotapi.FileBytes{
		Name:  "guide.pdf",
		Bytes: fileBytes,
	}

	msg := tgbotapi.NewDocumentUpload(userID, doc)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка при отправке PDF: %v", err)
	} else {
		log.Printf("📄 Файл отправлен пользователю %d", userID)
	}
}
