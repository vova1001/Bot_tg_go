package handlers

import (
	"CourseTg/config"
	"CourseTg/internal/payment"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var processedUpdates = make(map[int]bool)

func HandleUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if processedUpdates[update.UpdateID] {
		return
	}
	processedUpdates[update.UpdateID] = true

	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Всем привет 👋меня зовут Юлия, мне 39 лет.\n"+
				"Люблю силовые тренировки 💪🏼 и сбалансированное питание 🍽️\n"+
				"Общий стаж занятий 8 лет. Сертифицирована по направлению «Фитнес-занятия во время беременности и восстановления после родов».\n"+
				"Благодаря правильному и сбалансированному питанию добилась своей лучшей формы, что позволяет мне держать тело в тонусе.\n"+
				"В этом боте вы можете приобрести мои сборники рецептов простых и вкусных блюд, благодаря которым я получила подтянутую фигуру.")

			button := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📚 Показать продукты", "show_courses"),
				),
			)
			msg.ReplyMarkup = button
			bot.Send(msg)

		case "/exit":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот выключается. Но не переживайте, он продолжит работать.")
			bot.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я такое не понимаю, выберите другую команду.")
			bot.Send(msg)
		}
	}

	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data

		switch {
		case data == "show_courses":
			buttons := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📗 Сборник готовых завтраков", "course_2"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📘 Книга рецептов", "course_1"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📕 Книга рецептов + Сборник готовых завтраков", "course_3"),
				),
			)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите сборник для покупки:")
			msg.ReplyMarkup = buttons
			bot.Send(msg)

		case len(data) > 6 && data[:6] == "course":
			var courseDescription string
			switch data {
			case "course_2":
				courseDescription = "📗 Сборник готовых завтраков\n💰 Цена: 399₽"
			case "course_1":
				courseDescription = "📘 Книга рецептов\n💰 Цена: 599₽"
			case "course_3":
				courseDescription = "📕 Книга рецептов + Сборник готовых завтраков\n💰 Цена: 800₽"
			}

			buttons := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("💳 Забрать", "buy_"+data),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "show_courses"),
				),
			)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, courseDescription)
			msg.ReplyMarkup = buttons
			bot.Send(msg)

		case len(data) > 4 && data[:4] == "buy_":
			course := data[4:]
			var amount, desc string

			switch course {
			case "course_1":
				amount = "599.00"
				desc = "Книга рецептов"
			case "course_2":
				amount = "399.00"
				desc = "Сборник завтраков"
			case "course_3":
				amount = "800.00"
				desc = "Книга + Завтраки"
			default:
				return
			}

			cfg := config.LoadConfig()
			url, err := payment.CreatePayment(amount, desc, cfg.SuccessURL, cfg.YooKassaShopID, cfg.YooKassaSecretKey)
			if err != nil {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Ошибка при создании ссылки на оплату.")
				bot.Send(msg)
				return
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "💳 Перейдите по ссылке для оплаты:\n"+url)
			bot.Send(msg)
		}

		bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
	}
}
