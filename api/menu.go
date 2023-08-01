package api

import (
	"log"
	"os"
	"fmt"
	_ "github.com/TechSir3n/CityCompanion/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

//if update.Message.Location == nil {
//	bot.Send("Пожалуйста, поделитесь своими кординатами местоположения. Для поиска ближайщих выбранных вами мест")
//}

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

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			switch callbackData {
			case "action1":
				fmt.Println("Here action1")
			  // выполнение действий для кнопки "Парки-Атракционы"
			// и так далее для всех кнопок
			}
		  } else { 
			fmt.Println("Is nil")
		  }
		

		switch update.Message.Text {
		case "/start":
			reply := "Добро пожаловать в CityCompanion! Я ваш надежный гид по городу. " +
				"Просто отправьте мне свои координаты, и я помогу вам найти лучшие места в городе: от уютных кафе и ресторанов до кинотеатров и парков с аттракционами." +
				"Отправляйте свои запросы, и я с радостью помогу вам насладиться лучшими местами в вашем городе! "
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyMarkup = createMainMenu()

			bot.Send(msg)
		case "❗️Показать меню":
			msgN = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нужное действие: ")
			msgN.ReplyMarkup = createNeedAction()

			bot.Send(msgN)
		case "📍 Поделится с кординатами местоположения":
			setLocation(update)
		default:
			reply := "Я не понимаю, что вы говорите."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

			bot.Send(msg)
		}

	}
}

func createNeedAction() tgbotapi.InlineKeyboardMarkup {
	replyMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("🍽️ Кафе-Рестораны", "action1"),
				tgbotapi.NewInlineKeyboardButtonData("🎡🎢 Парки-Атракционы", "action2"),
				tgbotapi.NewInlineKeyboardButtonData("👨‍👩‍👧‍👦 Отдых с детьми", "action3"),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("🏝️ Пляжи", "action4"),
				tgbotapi.NewInlineKeyboardButtonData("🔍 Достопремечательности", "action5"),
				tgbotapi.NewInlineKeyboardButtonData("🎬 Просмотр фильмов", "action6"),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("🏋️‍♀️ Тренажерныe Залы", "action7"),
				tgbotapi.NewInlineKeyboardButtonData("🏃‍♀️ Cпорт площадки", "action8"),
				tgbotapi.NewInlineKeyboardButtonData("🎤 Караоке", "action9"),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("👩‍⚕️💉 Скорая помощь", "action10"),
				tgbotapi.NewInlineKeyboardButtonData("💵 Банкоматы", "action11"),
			},
		},
	}

	return replyMarkup
}

func createMainMenu() tgbotapi.ReplyKeyboardMarkup {
	menu := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("❗️Показать меню"),
				tgbotapi.NewKeyboardButton("📍 Поделится с кординатами местоположения"),
			},
		},
		ResizeKeyboard: true,
	}
	return menu
}
