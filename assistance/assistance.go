package assistance

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"

)

func WarningLocation(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	reply := "Перед тем как начать искать ближайщие интересующие вас места,нам следует знать ваши кординаты местоположения," +
		"только основываясь на них мы сможем предложить вам ближайщие и интереснейшие места."
	menu := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("📍 Поделится с кординатами местоположения"),
			},
		},
		ResizeKeyboard: true,
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyMarkup = menu
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("[WarningLocation] : %v", err)
	}
}

func AboutBot(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	reply := "CityCompanion - ваш личный помощник. Бот способен производить поиск самых ближайших и интересующих вас мест," +
		"а также находить заведения, исходя из отзывов и качества обслуживания в выбранных вами категориях"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("[AboutBot] : %v", err)
	}
}

func AskLimit(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) (int64, int64) {
	reply := "Пожалуйста, укажите максимальное количество мест, которое вы бы хотели видеть в списке:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5-2", "limit1"),
			tgbotapi.NewInlineKeyboardButtonData("7-2", "limit2"),
			tgbotapi.NewInlineKeyboardButtonData("10-2", "limit3"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("15-2", "limit4"),
			tgbotapi.NewInlineKeyboardButtonData("20-2", "limit5"),
			tgbotapi.NewInlineKeyboardButtonData("25-2", "limit6"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("[AskLimitPlaces] : %v", err)
	}

	var limitPlace int64
	for u := range updates {
		if u.CallbackQuery != nil {
			limitText := u.CallbackQuery.Data
			switch limitText {
			case "limit1":
				limitPlace = 5
			case "limit2":
				limitPlace = 7
			case "limit3":
				limitPlace = 10
			case "limit4":
				limitPlace = 15
			case "limit5":
				limitPlace = 20
			case "limit6":
				limitPlace = 25
			}
			break
		}
	}

	return 1, limitPlace
}

func AdjuctRadiusSearch(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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
}

func ParseRadius(radius string) float64 {
	if radius == "" {
		return 0.0
	}

	var meter float64
	if strings.HasSuffix(radius, "км") {
		radius = strings.TrimSuffix(radius, "км")
		num, _ := strconv.ParseFloat(radius, 64)
		meter = num * 1000

	} else if strings.HasSuffix(radius, "м") {
		radius = strings.TrimSuffix(radius, "м")
		meter, _ = strconv.ParseFloat(radius, 64)
	}
	return meter
}

