package api

import (
	"context"
	utils "github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
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
		utils.Error("Failed to save location user: %v", err.Error())
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваши кординаты успешно получены."+
		"Благодарим,теперь мы сможем находить самые ближайщие интересующие вас места.")
	bot.Send(msg)
}

func handleRadiusResponse(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	dbRadius := database.NewRadiusSearchImpl(database.DB)

	if strings.Contains(update.Message.Text, "Да") {
		reply := "Пожалуйста, введите радиус поиска, ограничивая его расстоянием в метрах или километрах,пример(1850м, 5км)"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		bot.Send(msg)

		var radius string
		for u := range updates {
			radius = u.Message.Text
			if utils.IsRadiusCorrect(radius) {
				meter := utils.ParseRadius(radius)
				err := dbRadius.SaveRadiusSearch(context.Background(), meter)
				if err != nil {
					utils.Error("Save radius error: ", err.Error())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Радиус успешно сохранен, и будет применён в поисках ближайщих мест")
				bot.Send(msg)
				break
			} else {
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Вы ввели неверный формат радиуса. Пожалуйста, повторите попытку.")
				bot.Send(msg)
			}
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
