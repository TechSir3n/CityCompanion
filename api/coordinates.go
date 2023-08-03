package api

import (
	"context"
	"math"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func AskCoordinates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Разрешите доступ к вашим кординатам местоположения ?")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
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
