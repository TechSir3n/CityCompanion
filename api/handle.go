package api

import (
	"context"
	"github.com/TechSir3n/CityCompanion/assistance"
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
	if err, lat, lon := userLocation.GetUserLocation(context.Background(), update.Message.Chat.ID); lat != 0.0 && lon != 0.0 {
		if err = userLocation.UpdateUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude); err != nil {
			utils.Error("Failed to save  update user: %v", err)
		}
	} else {
		if err := userLocation.SaveUserLocation(context.Background(), update.Message.Chat.ID, latitude, longitude); err != nil {
			utils.Error("Failed to save location user: %v", err)
		}

	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваши кординаты успешно получены."+
		"Благодарим,теперь мы сможем находить самые ближайщие интересующие вас места.")
	msg.ReplyMarkup = createMainMenu()
	bot.Send(msg)
}

func handleCordinates(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	textInput := make(chan string)
	waitInputUser(textInput, make(chan string), updates)
	text := <-textInput

	if text == "Запретить" {
		msgc := tgbotapi.NewMessage(update.Message.Chat.ID, "Для продолжения работы с ботом, пожалуйста, выберите орентир, город или адрес")

		keyboardc := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("\U0001F3E0 Город"),
				tgbotapi.NewKeyboardButton("\U0001F3E2 Улица"),
			),
		)

		keyboardc.OneTimeKeyboard = true
		keyboardc.ResizeKeyboard = true
		msgc.ReplyMarkup = keyboardc
		bot.Send(msgc)

		db := database.NewUserLocationImpl(database.DB)
		errCh := make(chan string)

		waitInputUser(textInput, errCh, updates)
		text := <-textInput
		if text == "\U0001F3E0 Город" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название города: ")
			bot.Send(msg)

			cityInput := make(chan string)
			waitInputUser(cityInput, errCh, updates)
			city := <-cityInput

			if latitude, longitude, err := getCordinatesByCity(city); err != nil ||
				longitude == 0.0 && latitude == 0.0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось определить кординаты вашего города,попробуйте ещё раз или"+
					"проверти корректность веденного вами города")
				bot.Send(msg)
				return
			} else {
				isHaveDB(bot, db, update, latitude, longitude)
			}
		} else if text == "\U0001F3E2 Улица" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название вашей улицы в формате(улица,город): ")
			bot.Send(msg)

			streetInput := make(chan string)
			waitInputUser(streetInput, errCh, updates)
			street := <-streetInput

			if !utils.IsRightFormat(street) {
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
				isHaveDB(bot, db, update, latitude, longitude)
			}
		}
	} else if text == "Разрешить \U0001F44D" {
		btn := tgbotapi.KeyboardButton{
			RequestLocation: true,
			Text:            "Подвердить",
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Благодарим ! Осталось подтвердить")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
		if _, err := bot.Send(msg); err != nil {
			utils.Fatal("Error sending message: %v", err)
		}

	}
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
			if _, prevRadius := dbRadius.GetRadiusSearch(context.Background(), update.Message.Chat.ID); prevRadius != 0.0 {
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

			msgM := tgbotapi.NewMessage(update.Message.Chat.ID, "Радиус успешно сохранен и будет применен в поисках ближайших мест.")
			msgM.ReplyMarkup = createMainMenu()
			bot.Send(msgM)
		} else {
			msg.Text = "Вы ввели неверный формат радиуса. Пожалуйста, повторите попытку."
			bot.Send(msg)
			handleRadiusResponse(bot, update, updates)
		}

	case "Нет":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Радиус поиска не будет ограничен.")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = createMainMenu()
		bot.Send(msg)
	default:
		break
	}
}

func handlePlaceCategory(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel, category string) {
	limitPhoto, limitPlace := assistance.AskLimit(bot, update, updates)
	if isCoordinatesShared(update.Message.Chat.ID) {
		GetNearbyPlaces(limitPlace, limitPhoto, category, bot, update, updates)
	} else {
		assistance.WarningLocation(bot, update)
	}
}
