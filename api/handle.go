package api

import (
	"context"
	utils "github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func handleGeocoding(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil || update.Message.Location == nil {
		return
	}

	location := update.Message.Location
	latitude := location.Latitude
	longitude := location.Longitude

	userLocation := database.NewUserLocationImpl(database.DB)
	err := userLocation.SaveUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude)
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

	switch update.Message.Text {
	case "Да":
		reply := "Пожалуйста, введите радиус поиска, ограничивая его расстоянием в метрах или километрах, пример (1850м, 5км)"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		bot.Send(msg)

		radiusInput := make(chan string)
		waitInputUser(radiusInput, make(chan string), updates)
		radius := <-radiusInput

		if utils.IsRadiusCorrect(radius) {
			meter := utils.ParseRadius(radius)
			_, prevRadius := dbRadius.GetRadiusSearch(context.Background(), update.Message.Chat.ID)
			if prevRadius != 0.0{
				if radiusInt, err := strconv.Atoi(radius); err == nil {
					radiusUp := float64(radiusInt)
					if err = dbRadius.UpdateRadiusSearch(context.Background(), update.Message.Chat.ID, radiusUp); err != nil {
						utils.Error(err)
					}
				}
			} else {
				if err := dbRadius.SaveRadiusSearch(context.Background(), update.Message.Chat.ID, meter); err != nil {
					utils.Error("Save radius error: ", err)
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Радиус успешно сохранен и будет применен в поисках ближайших мест.")
			bot.Send(msg)
			break
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели неверный формат радиуса. Пожалуйста, повторите попытку.")
			bot.Send(msg)
		}

	case "Нет":
		reply := "Радиус поиска не будет ограничен."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		msg.ReplyMarkup = createMainMenu()
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный ввод. Пожалуйста, используйте кнопки для выбора.")
		msg.ReplyMarkup = createMainMenu()
		bot.Send(msg)
	}
}
