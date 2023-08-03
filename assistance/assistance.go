package assistance

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func ParseRadius(radius string) {
	if radius == "" {
		return
	}

	if strings.HasSuffix(radius, "км") {
		radius = strings.TrimSuffix(radius, "км")
		num, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			Error("Failed to parse radius to float64", err.Error())
		} else {
			meter := num * 1000
			radius = strconv.Itoa(int(meter))
		}

	} else if strings.HasSuffix(radius, "м") {
		radius = strings.TrimSuffix(radius, "м")
	}

	fmt.Println("Here1: ",radius)
}

func IsRadiusCorrect(radius string) bool {
	if strings.HasSuffix(radius, "км") {
		radius = strings.TrimSuffix(radius, "км")
		num, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			return false
		} else if num < 0 {
			return false
		}
		return true
	} else if strings.HasSuffix(radius, "м") {
		radius = strings.TrimSuffix(radius, "м")
		num, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			return false
		} else if num < 0 {
			return false
		}
		return true
	}
	return false
}
