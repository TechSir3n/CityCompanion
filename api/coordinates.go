package api

import (
	"context"
	"fmt"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
)

func AskCoordinates(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	btnAllow := tgbotapi.KeyboardButton{Text: "Разрешить"}
	btnDeny := tgbotapi.KeyboardButton{RequestLocation: false, Text: "Запретить"}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(btnAllow),
		tgbotapi.NewKeyboardButtonRow(btnDeny),
	)

	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Разрешите доступ к вашим кординатам местоположения ?")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	handleCordinates(bot, update, updates)
}

func isHaveDB(bot *tgbotapi.BotAPI, db *database.UserLocationImpl, update tgbotapi.Update, latitude, longitude float64) {
	if err, lat, lon := db.GetUserLocation(context.Background(), update.Message.Chat.ID); err == nil {
		if lat != 0.0 && lon != 0.0 {
			if err = db.UpdateUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude); err != nil {
				fmt.Printf("Error update location: %v", err)
			}
		} else {
			if err = db.SaveUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude); err != nil {
				fmt.Printf("Error save location: %v", err)
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Благодарим,кординаты вашего города успешно получены, и сохранены !")
		bot.Send(msg)
	}
}

func isCoordinatesShared(userID int64) bool {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background(), userID)
	if err != nil {
		return false
	} else if math.IsNaN(latitude) || math.IsNaN(longitude) {
		return false
	}

	return true
}
