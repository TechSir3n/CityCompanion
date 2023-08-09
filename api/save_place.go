package api

import (
	"context"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func savePlace(bot *tgbotapi.BotAPI, chatID int64, placeName string, locations []Location) {
	found := false
	db := database.NewSavedPlacesImpl(database.DB)
	for _, location := range locations {
		if placeName == location.Name {
			db.SavePlace(context.Background(), location.Name, location.Location.Address)
			msg := tgbotapi.NewMessage(chatID, "Отлично,место успешно сохранено")
			bot.Send(msg)
			found = true
			return
		}
	}

	if !found {
		msg := tgbotapi.NewMessage(chatID, "Простите,но такого мест нет в спискe.")
		bot.Send(msg)
	}
}

func saveFavoritePlace(bot *tgbotapi.BotAPI, chatID int64, placeName string, locations []Location) {
	found := false
	db := database.NewFavoritePlacesImp(database.DB)
	for _, location := range locations {
		if placeName == location.Name {
			db.SaveFavoritePlace(context.Background(), location.Name, location.Location.Address)
			msg := tgbotapi.NewMessage(chatID, "Отлично,место успешно добавлено в избранное")
			bot.Send(msg)
			found = true
			return
		}
	}

	if !found {
		msg := tgbotapi.NewMessage(chatID, "Простите,но такого мест нет в спискe.")
		bot.Send(msg)
	}
}
