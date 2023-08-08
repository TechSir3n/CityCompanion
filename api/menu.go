package api

import (
	"context"
	"fmt"
	"github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
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
	fmt.Printf("Authorized on account %s", bot.Self.UserName)

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
		case "/getmylocation":
			street := GetUserStreet()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, street)
			bot.Send(msg)
		case "/favoriteplace":
			f_db := database.NewFavoritePlacesImp(database.DB)
			if names, addresses, err := f_db.GetFavoritePlaces(context.Background()); err != nil {
				assistance.Error(err.Error())
			} else {
				var message string
				for i := range names {
					message += fmt.Sprintf("%d. %s - %s\n", i+1, names[i], addresses[i])
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				bot.Send(msg)
			}
		case "/savedplaces":
			s_db := database.NewSavedPlacesImpl(database.DB)
			if names, addresses, err := s_db.GetSavePlaces(context.Background()); err != nil {
				assistance.Error(err.Error())
			} else {
				var message string
				for i := range names {
					message += fmt.Sprintf("%d. %s - %s\n", i+1, names[i], addresses[i])
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				bot.Send(msg)
			}
		case "/adjustradius":
			assistance.AdjuctRadiusSearch(bot, update)
		case "🍽️ Кафе-Рестораны":
			handlePlaceCategory(bot, update, updates, "13000")
		case "🍵 Кофейная-Чайная":
			handlePlaceCategory(bot, update, updates, "13035")
		case "🍣 Японская кухня":
			handlePlaceCategory(bot, update, updates, "13276")
		case "🏨 Отели":
			handlePlaceCategory(bot, update, updates, "19014")
		case "🍰 Кондитерские магазины":
			handlePlaceCategory(bot, update, updates, "17057")
		case "🏖️ Пляжи":
			handlePlaceCategory(bot, update, updates, "16003")
		case "🏛️ Достопремечательности и природа":
			handlePlaceCategory(bot, update, updates, "16020")
		case "🌳 Городские парки":
			handlePlaceCategory(bot, update, updates, "16032")
		case "🏋️‍♀️ Тренажерный зал и студии":
			handlePlaceCategory(bot, update, updates, "19066")
		case "💆‍♀️ Услуги для здоровья и красоты":
			handlePlaceCategory(bot, update, updates, "17035")
		case "⛪️ Церквки-Мечети":
			handlePlaceCategory(bot, update, updates, "12106")
		case "🛍️ Магазины":
			handlePlaceCategory(bot, update, updates, "17096")
		case "🍻 Бары":
			handlePlaceCategory(bot, update, updates, "13012")
		default:
			handleRadiusResponse(bot, update, updates)
			handleGeocoding(bot, update)
		}
	}
}

func handlePlaceCategory(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel, category string) {
	limitPhoto, limitPlace := assistance.AskLimit(bot, update, updates)
	if isCoordinatesShared() {
		GetNearbyPlaces(limitPlace, limitPhoto, category, bot, update, updates)
	} else {
		assistance.WarningLocation(bot, update)
	}
}

func createNeedAction() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("🍽️ Кафе-Рестораны"),
				tgbotapi.NewKeyboardButton("🍵 Кофейная-Чайная"),
				tgbotapi.NewKeyboardButton("🍣 Японская кухня"),
			},
			{
				tgbotapi.NewKeyboardButton("🏖️ Пляжи"),
				tgbotapi.NewKeyboardButton("🏛️ Достопремечательности и природа"),
				tgbotapi.NewKeyboardButton("🌳 Городские парки"),
			},
			{
				tgbotapi.NewKeyboardButton("🏋️‍♀️ Тренажерный зал и студии"),
				tgbotapi.NewKeyboardButton("💆‍♀️ Услуги для здоровья и красоты"),
				tgbotapi.NewKeyboardButton("⛪️ Церквки-Мечети"),
			},
			{
				tgbotapi.NewKeyboardButton("🛍️ Магазины"),
				tgbotapi.NewKeyboardButton("🍻 Бары"),
				tgbotapi.NewKeyboardButton("🍰 Кондитерские магазины"),
				tgbotapi.NewKeyboardButton("🏨 Отели"),
			},
		},
		ResizeKeyboard: true,
	}
	replyMarkup.OneTimeKeyboard = true

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
