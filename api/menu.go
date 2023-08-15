package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	_ "github.com/TechSir3n/CityCompanion/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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

		commands := assistance.NewComnands()
		categories := assistance.NewPlaceCategories()
		categoriesCode := assistance.NewCodeCategories()

		switch update.Message.Text {
		case commands.Start:
			reply := "Добро пожаловать в CityCompanion! Я ваш надежный гид по городу. " +
				"Просто отправьте мне свои координаты, и я помогу вам найти лучшие места в городе: от уютных кафе и ресторанов до кинотеатров и парков с аттракционами." +
				"Отправляйте свои запросы, и я с радостью помогу вам насладиться лучшими местами в вашем городе! "
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyMarkup = createMainMenu()
			bot.Send(msg)
		case commands.Menu:
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()
			bot.Send(msgN)
		case commands.Coordinates:
			AskCoordinates(bot, update, updates)
		case commands.Radius:
			assistance.AdjuctRadiusSearch(bot, update)
		case commands.About:
			assistance.AboutBot(bot, update)
		case commands.Menu:
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()
			bot.Send(msgN)
		case commands.SendLocation:
			AskCoordinates(bot, update, updates)
		case commands.GetLocation:
			street := GetUserStreet(update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, street)
			bot.Send(msg)
		case commands.FavoritePlace:
			f_db := database.NewFavoritePlacesImp(database.DB)
			if names, addresses, err := f_db.GetFavoritePlaces(context.Background(), update.Message.Chat.ID); err != nil {
				assistance.Error(err.Error())
			} else {
				var message string
				for i := range names {
					message += fmt.Sprintf("%d. %s - %s\n", i+1, names[i], addresses[i])
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				bot.Send(msg)
			}
		case commands.SavedPlace:
			s_db := database.NewSavedPlacesImpl(database.DB)
			if names, addresses, err := s_db.GetSavePlaces(context.Background(), update.Message.Chat.ID); err != nil {
				assistance.Error(err.Error())
			} else {
				var message string
				for i := range names {
					message += fmt.Sprintf("%d. %s - %s\n", i+1, names[i], addresses[i])
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				bot.Send(msg)
			}
		case commands.AdjustRadius:
			assistance.AdjuctRadiusSearch(bot, update)
		case categories.CafeAndRestaurants:
			handlePlaceCategory(bot, update, updates, categoriesCode.CafeAndRestaurantsCode)
		case categories.CoffeeAndTea:
			handlePlaceCategory(bot, update, updates, categoriesCode.CoffeeAndTeaCode)
		case categories.JapaneseFood:
			handlePlaceCategory(bot, update, updates, categoriesCode.JapaneseFoodCode)
		case categories.Hotels:
			handlePlaceCategory(bot, update, updates, categoriesCode.HotelsCode)
		case categories.ConfectioneryStores:
			handlePlaceCategory(bot, update, updates, categoriesCode.ConfectioneryStoresCode)
		case categories.Beaches:
			handlePlaceCategory(bot, update, updates, categoriesCode.BeachesCode)
		case categories.SightsAndNature:
			handlePlaceCategory(bot, update, updates, categoriesCode.SightsAndNatureCode)
		case categories.CityParks:
			handlePlaceCategory(bot, update, updates, categoriesCode.CityParksCode)
		case categories.GymAndStudios:
			handlePlaceCategory(bot, update, updates, categoriesCode.GymAndStudiosCode)
		case categories.HealthAndBeautyServices:
			handlePlaceCategory(bot, update, updates, categoriesCode.HealthAndBeautyServicesCode)
		case categories.ChurchesAndMosques:
			handlePlaceCategory(bot, update, updates, categoriesCode.ChurchesAndMosquesCode)
		case categories.Shops:
			handlePlaceCategory(bot, update, updates, categoriesCode.ShopsCode)
		case categories.Bars:
			handlePlaceCategory(bot, update, updates, categoriesCode.BarsCode)
		default:
			handleRadiusResponse(bot, update, updates)
			handleGeocoding(bot, update)
		}
	}
}

func createNeedAction() tgbotapi.ReplyKeyboardMarkup {
	categories := assistance.NewPlaceCategories()
	replyMarkup := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton(categories.CafeAndRestaurants),
				tgbotapi.NewKeyboardButton(categories.CoffeeAndTea),
				tgbotapi.NewKeyboardButton(categories.JapaneseFood),
			},
			{
				tgbotapi.NewKeyboardButton(categories.Beaches),
				tgbotapi.NewKeyboardButton(categories.SightsAndNature),
				tgbotapi.NewKeyboardButton(categories.CityParks),
			},
			{
				tgbotapi.NewKeyboardButton(categories.GymAndStudios),
				tgbotapi.NewKeyboardButton(categories.HealthAndBeautyServices),
				tgbotapi.NewKeyboardButton(categories.ChurchesAndMosques),
			},
			{
				tgbotapi.NewKeyboardButton(categories.Shops),
				tgbotapi.NewKeyboardButton(categories.Bars),
				tgbotapi.NewKeyboardButton(categories.ConfectioneryStores),
				tgbotapi.NewKeyboardButton(categories.Hotels),
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
