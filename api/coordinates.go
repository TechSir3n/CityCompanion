package api

import (
	"context"
	"math"

	"github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

		db := database.NewUserLocationImpl(database.DB)

		getNewUpdate := <-updates
		if getNewUpdate.Message != nil && getNewUpdate.Message.Text == "Город" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название города: ")
			bot.Send(msg)

			city := make(chan string)
			waitInputUser(city, updates)

			if latitude, longitude, err := getCordinatesByCity(<-city); err != nil ||
				longitude == 0.0 && latitude == 0.0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось определить кординаты вашего города,попробуйте ещё раз или"+
					"проверти корректность веденного вами города")
				bot.Send(msg)
				return
			} else {
				db.SaveUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Благодарим,кординаты вашего города успешно получены, и сохранены !")
				bot.Send(msg)
			}
		} else if getNewUpdate.Message.Text == "Улица" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название вашей улицы в формате(улица,город): ")
			bot.Send(msg)

			streetInput := make(chan string)
			waitInputUser(streetInput, updates)
			street := <-streetInput

			if !assistance.IsRightFormat(street) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели неверный формат улицы")
				bot.Send(msg)
				return
			}

			if latitude, longitude, err := getCoordinatesByStreet(street); err != nil ||
				latitude == 0.0 && longitude == 0.0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось определить кординаты вашей улицы,попробуйте ещё раз или "+
					"проверти корректность веденного вами улицы")
				bot.Send(msg)
			} else {
				db.SaveUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Благодарим,кординаты вашей улицы успешно получены, и сохранены !")
				bot.Send(msg)
			}

		}
	}
}

func isCoordinatesShared(userID int64) bool {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background(),userID)
	if err != nil {
		return false
	} else if math.IsNaN(latitude) || math.IsNaN(longitude) {
		return false
	}

	return true
}
