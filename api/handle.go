package api

import (
	"context"
	"fmt"
	"log"
	"strings"

	utils "github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// we get here latitude and longitude user's
func handleGeocoding(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil || update.Message.Location == nil {
		return
	}

	location := update.Message.Location
	latitude := location.Latitude
	longitude := location.Longitude

	userLocation := database.NewUserLocationImpl(database.DB)
	err := userLocation.SaveUserLocation(context.Background(), latitude, longitude)
	if err != nil {
		log.Fatalf("Failed save user location: %v", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваши кординаты успешно получены."+
		"Благодарим,теперь мы сможем находить самые ближайщие интересующие вас места.")
	bot.Send(msg)
}

func handleRadiusResponse(bot *tgbotapi.BotAPI, update tgbotapi.Update, updConfig tgbotapi.UpdateConfig) {
	if strings.Contains(update.Message.Text, "Да") {
		reply := "Пожалуйста, введите радиус поиска, ограничивая его расстоянием в метрах или километрах,пример(1850м, 5км)"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		bot.Send(msg)

		updates, err := bot.GetUpdates(updConfig)
		if err != nil {
			utils.Error("Failed to get updates", err.Error())
			return
		}

		var radius string
		for _, upd := range updates {
			if upd.Message != nil && upd.Message.Chat.ID == update.Message.Chat.ID {
				radius = upd.Message.Text
				break
			}
		}

		fmt.Println("Radius: ", radius)

		if utils.IsRadiusCorrect(radius) {
			utils.ParseRadius(radius)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели неверный формат радиуса")
			bot.Send(msg)
		}
	} else if update.Message.Text == "Нет" {
		reply := "Радиус поиска не будет ограничен"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		bot.Send(msg)
	} else {
		reply := "Пожалуйста выберите один из вариантов ответа"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		yesBTN := tgbotapi.NewKeyboardButton("Да")
		noBTN := tgbotapi.NewKeyboardButton("Нет")

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(yesBTN, noBTN),
		)

		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}
