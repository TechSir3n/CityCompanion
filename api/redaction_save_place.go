package api

import (
	"context"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func redactionSavePlace(bot *tgbotapi.BotAPI, upd tgbotapi.Update, chatID int64, updates tgbotapi.UpdatesChannel) {
	savedPlacesDB := database.NewSavedPlacesImpl(database.DB)
	if upd.CallbackQuery.Data == "buttonDelete" {
		msg := tgbotapi.NewMessage(chatID, "Введите название места которое хотите удалить: ")
		bot.Send(msg)

		placeInput := make(chan string)
		waitInputUser(placeInput, nil, updates)
		placeName := <-placeInput

		if placeName != "" {
			if err := savedPlacesDB.DeleteOnePlace(context.Background(), chatID, placeName); err == nil {
				msg.Text = "Место успешно удалённо !"
				bot.Send(msg)
			} else {
				msg.Text = "Не удалось удалить место, вероятно такого места нет в списке!"
				bot.Send(msg)
				redactionSavePlace(bot, upd, chatID, updates)
			}
		}
	} else if upd.CallbackQuery.Data == "buttonClean" {
		if err := savedPlacesDB.DeleteSavePlaces(context.Background(), chatID); err == nil {
			msg := tgbotapi.NewMessage(chatID, "Отлично,список сохранённых мест очистен!")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Не удалось удалить место, вероятно такого места нет в списке!")
			bot.Send(msg)
			redactionSavePlace(bot, upd, chatID, updates)
		}
	} else {
		return
	}
}

func redactionFavoritePlace(bot *tgbotapi.BotAPI, upd tgbotapi.Update, chatID int64, updates tgbotapi.UpdatesChannel) {
	favoritePlacesDB := database.NewFavoritePlacesImp(database.DB)
	if upd.CallbackQuery.Data == "buttonDelete" {
		msg := tgbotapi.NewMessage(chatID, "Введите название места которое хотите удалить: ")
		bot.Send(msg)

		placeInput := make(chan string)
		waitInputUser(placeInput, nil, updates)
		placeName := <-placeInput

		if placeName != "" {
			if err := favoritePlacesDB.DeleteOneFavoritePlace(context.Background(), chatID, placeName); err == nil {
				msg.Text = "Место успешно удалённо !"
				bot.Send(msg)
			} else {
				msg.Text = "Не удалось удалить место, вероятно такого места нет в списке!"
				bot.Send(msg)
				redactionFavoritePlace(bot, upd, chatID, updates)
			}
		}
	} else if upd.CallbackQuery.Data == "buttonClean" {
		if err := favoritePlacesDB.DeleteFavoritePlaces(context.Background(), chatID); err == nil {
			msg := tgbotapi.NewMessage(chatID, "Отлично,список избранных мест очистен!")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Не удалось удалить место, вероятно такого места нет в списке!")
			bot.Send(msg)
			redactionFavoritePlace(bot, upd, chatID, updates)
		}
	} else {
		return
	}
}
