package api

import (
	"context"
	"fmt"
	"github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandsBot(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	commands := assistance.NewComnands()
	var msgN tgbotapi.MessageConfig
	var msg tgbotapi.MessageConfig

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("\U0001F5D1  Удалить", "buttonDelete"),
			tgbotapi.NewInlineKeyboardButtonData("\U0001F4A6 Очистить", "buttonClean"),
		),
	)

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
		if street, err := getUserStreet(update.Message.Chat.ID); street == "" || err != nil {
			fmt.Printf("User street error : %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, поделитесь своим местоположением.")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, street)
			bot.Send(msg)
		}
	case commands.FavoritePlace:
		f_db := database.NewFavoritePlacesImp(database.DB)
		if names, addresses, err := f_db.GetFavoritePlaces(context.Background(), update.Message.Chat.ID); err != nil ||
			names == nil && addresses == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "В вашем списке еще нёт избранных мест")
			bot.Send(msg)
		} else {
			var message string
			for i := range names {
				message += fmt.Sprintf("%d. %s - %s\n", i+1, names[i], addresses[i])
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ReplyMarkup = inlineKeyboard
			bot.Send(msg)

			for upd := range updates {
				if upd.CallbackQuery != nil {
					redactionFavoritePlace(bot, update, update.Message.Chat.ID, updates)
					break
				} else {
					return
				}
			}
		}

	case commands.SavedPlace:
		s_db := database.NewSavedPlacesImpl(database.DB)
		if names, addresses, err := s_db.GetSavePlaces(context.Background(), update.Message.Chat.ID); err != nil ||
			names == nil && addresses == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "В вашем списке еще нёт сохранённых мест")
			bot.Send(msg)
		} else {
			var message string
			for i := range names {
				message += fmt.Sprintf("%d. %s - %s\n", i+1, names[i], addresses[i])
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ReplyMarkup = inlineKeyboard
			bot.Send(msg)

			for upd := range updates {
				if upd.CallbackQuery != nil {
					redactionSavePlace(bot, upd, update.Message.Chat.ID, updates)
					break
				} else {
					return
				}
			}

		}

	case commands.AdjustRadius:
		assistance.AdjuctRadiusSearch(bot, update)

	default:
		handleRadiusResponse(bot, update, updates)
		handleGeocoding(bot, update)
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
