package api

import (
	"github.com/TechSir3n/CityCompanion/assistance"
	_ "github.com/TechSir3n/CityCompanion/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

func CreateButton() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error[CreateButtonFunc] : %v", err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	var msgN tgbotapi.MessageConfig
	var msg tgbotapi.MessageConfig

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			reply := "Добро пожаловать в CityCompanion! Я ваш надежный гид по городу. " +
				"Просто отправьте мне свои координаты, и я помогу вам найти лучшие места в городе: от уютных кафе и ресторанов до кинотеатров и парков с аттракционами." +
				"Отправляйте свои запросы, и я с радостью помогу вам насладиться лучшими местами в вашем городе! "
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyMarkup = createMainMenu()

			bot.Send(msg)
			break
		case "❗️Показать меню":
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()
			bot.Send(msgN)
			break
		case "📍 Поделится с кординатами местоположения":
			AskCoordinates(bot, update)
			break
		case "🔍 Настроить радиус поиска":
			reply := "Желаете ли вы ограничить радиус поиска интересующих вас мест?" +
				" Это позволить боту искать места, не превыщающие заданный радиус(расстояние),таким образом бот будет искать максимально приблежённые места от места вашего пребывания"

			yesBTN := tgbotapi.NewKeyboardButton("Да")
			noBTN := tgbotapi.NewKeyboardButton("Нет")

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(yesBTN, noBTN),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			break
		case "/about":
			assistance.AboutBot(bot, update)
			break
		case "/showmenu":
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()
			bot.Send(msgN)

			break
		case "/sendlocation":
			AskCoordinates(bot, update)
			break
		case "/adjustradius":
			break

		case "🍽️ Кафе-Рестораны":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🎡🎢 Парки-Атракционы":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "👨‍👩‍👧‍👦 Отдых с детьми":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🏝️ Пляжи":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🔍 Достопремечательности":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🎬 Просмотр фильмов":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🏋️‍♀️ Тренажерныe Залы":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🏃‍♀️ Cпорт площадки":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "🎤 Караоке":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "👩‍⚕️💉 Скорая помощь":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		case "💵 Банкоматы":
			if isCoordinatesShared() {

			} else {
				assistance.WarningLocation(bot, update)
			}

			break

		default:
			handleRadiusResponse(bot, update,updates)
			handleGeocoding(bot, update)
			break
		}

	}
}

func createNeedAction() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("🍽️ Кафе-Рестораны"),
				tgbotapi.NewKeyboardButton("🎡🎢 Парки-Атракционы"),
				tgbotapi.NewKeyboardButton("👨‍👩‍👧‍👦 Отдых с детьми"),
			},
			{
				tgbotapi.NewKeyboardButton("🏝️ Пляжи"),
				tgbotapi.NewKeyboardButton("🔍 Достопремечательности"),
				tgbotapi.NewKeyboardButton("🎬 Просмотр фильмов"),
			},
			{
				tgbotapi.NewKeyboardButton("🏋️‍♀️ Тренажерныe Залы"),
				tgbotapi.NewKeyboardButton("🏃‍♀️ Cпорт площадки"),
				tgbotapi.NewKeyboardButton("🎤 Караоке"),
			},
			{
				tgbotapi.NewKeyboardButton("👩‍⚕️💉 Скорая помощь"),
				tgbotapi.NewKeyboardButton("💵 Банкоматы"),
			},
		},
		ResizeKeyboard: true,
	}

	return replyMarkup
}

func createMainMenu() tgbotapi.ReplyKeyboardMarkup {
	menu := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("❗️Показать меню"),
				tgbotapi.NewKeyboardButton("📍 Поделится с кординатами местоположения"),
				tgbotapi.NewKeyboardButton("🔍 Настроить радиус поиска"),
			},
		},
		ResizeKeyboard: true,
	}
	return menu
}
