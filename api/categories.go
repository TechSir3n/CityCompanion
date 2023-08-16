package api

import (
	"github.com/TechSir3n/CityCompanion/assistance"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func categoriesPlace(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	categories := assistance.NewPlaceCategories()
	categoriesCode := assistance.NewPlaceCategoriesCode()

	switch update.Message.Text {
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
	}
}
