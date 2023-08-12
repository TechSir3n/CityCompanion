package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func showInMap(place  string,bot *tgbotapi.BotAPI, chatID int64, locations []Location) {
	found := false
	for _, location := range locations {
		if place == location.Name {
			msg := tgbotapi.NewLocation(chatID,
				location.Geocodes.Main.Latitude, location.Geocodes.Main.Longitude)
			if _, err := bot.Send(msg); err != nil {
				return
			}
			found = true
			return
		}
	}

	if !found {
		msg := tgbotapi.NewMessage(chatID, "Простите,но такого мест нет в спискe.")
		if _, err := bot.Send(msg); err != nil {
			return
		}
	}
}
