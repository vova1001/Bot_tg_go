package handlers

import (
	"CourseTg/config"
	"CourseTg/internal/payment"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserState struct {
	Step           string
	SelectedCourse string
	Amount         string
	Description    string
	Name           string
	Email          string
}

var processedUpdates = make(map[int]bool)
var UserStates = make(map[int64]*UserState)

func HandleUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if processedUpdates[update.UpdateID] {
		return
	}
	processedUpdates[update.UpdateID] = true

	if update.Message != nil {
		chatID := update.Message.Chat.ID

		state, exists := UserStates[chatID]
		if exists {
			switch state.Step {
			case "wait_name":
				state.Name = update.Message.Text
				state.Step = "wait_email"
				msg := tgbotapi.NewMessage(chatID, "✉️ Введите ваш email для получения чека:")
				bot.Send(msg)
				return
			case "wait_email":
				state.Email = update.Message.Text
				state.Step = "done"

				telegramID := fmt.Sprint(update.Message.From.ID)
				cfg := config.LoadConfig()
				url, err := payment.CreatePayment(
					state.Amount,
					state.Description,
					telegramID,
					state.SelectedCourse,
					state.Name,
					state.Email,
					cfg.YooKassaShopID,
					cfg.YooKassaSecretKey,
				)

				if err != nil {
					msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при создании ссылки на оплату.")
					bot.Send(msg)
					return
				}

				msg := tgbotapi.NewMessage(chatID, "💳 Перейдите по ссылке для оплаты:\n"+url)
				bot.Send(msg)
				return
			}
		}

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Всем привет 👋 меня зовут Юлия, мне 39 лет.\n\n"+
				"Люблю силовые тренировки 💪🏼 и сбалансированное питание 🍽️\n\n"+
				"Общий стаж занятий — 8 лет. Сертифицирована по направлению «Фитнес-занятия во время беременности и восстановления после родов».\n\n"+
				"Благодаря правильному и сбалансированному питанию добилась своей лучшей формы, что позволяет мне держать тело в тонусе.\n\n"+
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
		chatID := update.CallbackQuery.Message.Chat.ID
		data := update.CallbackQuery.Data

		switch data {
		case "show_courses":
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
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📋 Меню рационов", "show_menu"),
				),
			)
			msg := tgbotapi.NewMessage(chatID, "Выберите сборник для покупки или меню рационов:")
			msg.ReplyMarkup = buttons
			bot.Send(msg)

		case "show_menu":
			buttons := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🥗 Рацион на 1300 калорий", "ration_1"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🍱 Рацион на 1400 калорий", "ration_2"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🍛 Рацион на 1500 калорий", "ration_3"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🥙 Рацион на 1600 калорий", "ration_4"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🍜 Рацион на 1700 калорий", "ration_5"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🍲 Рацион на 1800 калорий", "ration_6"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "show_courses"),
				),
			)
			msg := tgbotapi.NewMessage(chatID, "Выберите рацион:")
			msg.ReplyMarkup = buttons
			bot.Send(msg)

		case "course_1", "course_2", "course_3",
			"ration_1", "ration_2", "ration_3", "ration_4", "ration_5", "ration_6":

			var courseDescription string
			switch data {
			case "course_2":
				courseDescription = "📗 Сборник готовых завтраков\n💰 Цена: 399₽"
			case "course_1":
				courseDescription = "📘 Книга рецептов\n💰 Цена: 599₽"
			case "course_3":
				courseDescription = "📕 Книга рецептов + Сборник готовых завтраков\n💰 Цена: 799₽"

			case "ration_1":
				courseDescription = "🥗 Рацион на 1300 калорий\n💰 Цена: 1000₽"
			case "ration_2":
				courseDescription = "🍱 Рацион на 1400 калорий\n💰 Цена: 1000₽"
			case "ration_3":
				courseDescription = "🍛 Рацион на 1500 калорий\n💰 Цена: 1000₽"
			case "ration_4":
				courseDescription = "🥙 Рацион на 1600 калорий\n💰 Цена: 1000₽"
			case "ration_5":
				courseDescription = "🍜 Рацион на 1700 калорий\n💰 Цена: 1000₽"
			case "ration_6":
				courseDescription = "🍲 Рацион на 1800 калорий\n💰 Цена: 1000₽"
			}

			buttons := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("💳 Забрать", "buy_"+data),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "show_menu"),
				),
			)
			msg := tgbotapi.NewMessage(chatID, courseDescription)
			msg.ReplyMarkup = buttons
			bot.Send(msg)

		default:
			if len(data) > 4 && data[:4] == "buy_" {
				course := data[4:]
				var amount, desc string

				switch course {
				case "course_2":
					amount = "1.00"
					desc = "Сборник готовых завтраков"
				case "course_1":
					amount = "1.00"
					desc = "Книга рецептов"
				case "course_3":
					amount = "1.00"
					desc = "Книга рецептов + Сборник готовых завтраков"

				case "ration_1":
					amount = "1.00"
					desc = "Рацион на 1300 калорий"
				case "ration_2":
					amount = "1.00"
					desc = "Рацион на 1400 калорий"
				case "ration_3":
					amount = "1.00"
					desc = "Рацион на 1500 калорий"
				case "ration_4":
					amount = "1.00"
					desc = "Рацион на 1600 калорий"
				case "ration_5":
					amount = "1.00"
					desc = "Рацион на 1700 калорий"
				case "ration_6":
					amount = "1.00"
					desc = "Рацион на 1800 калорий"
				default:
					return
				}

				UserStates[chatID] = &UserState{
					Step:           "wait_name",
					SelectedCourse: course,
					Amount:         amount,
					Description:    desc,
				}

				msg := tgbotapi.NewMessage(chatID, "🧾 Введите ваше имя для оформления чека:")
				bot.Send(msg)
			}
		}
	}
}
