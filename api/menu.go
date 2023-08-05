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
		case "❗️Показать меню":
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()
			bot.Send(msgN)
		case "📍 Поделится с кординатами местоположения":
			AskCoordinates(bot, update)
		case "🔍 Настроить радиус поиска":
			assistance.AdjuctRadiusSearch(bot, update)
		case "/about":
			assistance.AboutBot(bot, update)
		case "/showmenu":
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()
			bot.Send(msgN)
		case "/sendlocation":
			AskCoordinates(bot, update)
		case "/adjustradius":
			assistance.AdjuctRadiusSearch(bot, update)
		case "🍽️ Рестораны":
			handlePlaceCategory(bot, update, updates, "13065")
		case "🍵 Кафе, Кофейни и Чайные Дома":
			handlePlaceCategory(bot, update, updates, "13032")
		case "🛒 Розничная торговля продуктами питания и напитками":
			handlePlaceCategory(bot, update, updates, "17142")
		case "🏖️ Пляжи":
			handlePlaceCategory(bot, update, updates, "16003")
		case "🏛️ Достопремечательности":
			handlePlaceCategory(bot, update, updates, "16000")
		case "🌳 Городские парки":
			handlePlaceCategory(bot, update, updates, "16032")
		case "🏋️‍♀️ Тренажерный зал и студии":
			handlePlaceCategory(bot, update, updates, "18021")
		case "💆‍♀️ Услуги для здоровья и красоты":
			handlePlaceCategory(bot, update, updates, "11061")
		case "💇‍♂️ Парикмахерские":
			handlePlaceCategory(bot, update, updates, "11062")
		case "🛍️ Магазины одежды":
			handlePlaceCategory(bot, update, updates, "17043")
		case "🍻 Бары":
			handlePlaceCategory(bot, update, updates, "13003")
		default:
			handleRadiusResponse(bot, update, updates)
			handleGeocoding(bot, update)
		}
	}
}

func handlePlaceCategory(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel, category string) {
	limitPhoto, limitPlace := assistance.AskLimit(bot, update, updates)
	if isCoordinatesShared() {
		GetNearbyPlaces(limitPlace, limitPhoto, category, bot, update)
	} else {
		assistance.WarningLocation(bot, update)
	}
}

func createNeedAction() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("🍽️ Рестораны"),
				tgbotapi.NewKeyboardButton("🍵 Кафе, Кофейни и Чайные Дома"),
				tgbotapi.NewKeyboardButton("🛒 Розничная торговля продуктами питания и напитками"),
			},
			{
				tgbotapi.NewKeyboardButton("🏖️ Пляжи"),
				tgbotapi.NewKeyboardButton("🏛️ Достопремечательности"),
				tgbotapi.NewKeyboardButton("🌳 Городские парки"),
			},
			{
				tgbotapi.NewKeyboardButton("🏋️‍♀️ Тренажерный зал и студии"),
				tgbotapi.NewKeyboardButton("💆‍♀️ Услуги для здоровья и красоты"),
				tgbotapi.NewKeyboardButton("💇‍♂️ Парикмахерские"),
			},
			{
				tgbotapi.NewKeyboardButtonContact("🛍️ Магазины одежды"),
				tgbotapi.NewKeyboardButton("🍻 Бары"),
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
