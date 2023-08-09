package api

import (
	"context"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
)

func AskCoordinates(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	btnAllow := tgbotapi.KeyboardButton{
		RequestLocation: true,
		Text:            "Разрешить",
	}

	btnDeny := tgbotapi.KeyboardButton{
		RequestLocation: false,
		Text:            "Запретить",
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(btnAllow),
		tgbotapi.NewKeyboardButtonRow(btnDeny),
	)

	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Разрешите доступ к вашим кординатам местоположения ?")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	getNewUpdate := <-updates
	if getNewUpdate.Message.Text != "" && getNewUpdate.Message.Text == "Запретить" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для продолжения работы с ботом, пожалуйста, выберите орентир, город или адрес")

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Город"),
				tgbotapi.NewKeyboardButton("Улица"),
			),
		)

		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true
		msg.ReplyMarkup = keyboard
		bot.Send(msg)

		getNewUpdate := <-updates
		if getNewUpdate.Message != nil && getNewUpdate.Message.Text == "Город" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название города: ")
			bot.Send(msg)

			city := make(chan string)
			waitInputUser(city, updates)
			db := database.NewUserLocationImpl(database.DB)

			if latitude, longitude, err := GetCordinatesByCity(<-city); err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось определить кординаты вашего города,попробуйте ещё раз или"+
					"проверти корректность веденного вами города")
				bot.Send(msg)
			} else {
				db.SaveUserLocation(context.Background(), latitude, longitude)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Благодарим,кординаты вашего города успешно получены, и сохранены !")
				bot.Send(msg)
			}
		} else if getNewUpdate.Message.Text == "Улица" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название вашей улицы: ")
			bot.Send(msg)

			street := make(chan string)
			waitInputUser(street, updates)

		}
	}
}

func isCoordinatesShared() bool {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background())
	if err != nil {
		return false
	} else if math.IsNaN(latitude) || math.IsNaN(longitude) {
		return false
	}

	return true
}
