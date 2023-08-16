package assistance

type Commands struct {
	Start string

	About string

	Menu string

	Coordinates string

	Radius string

	SendLocation string

	GetLocation string

	FavoritePlace string

	SavedPlace string

	AdjustRadius string
}

func NewComnands() *Commands {
	return &Commands{
		Start:         "/start",
		About:         "/about",
		SendLocation:  "/sendlocation",
		GetLocation:   "/getmylocation",
		FavoritePlace: "/favoriteplace",
		SavedPlace:    "/savedplaces",
		AdjustRadius:  "/adjustradius",
		Menu:          "❗️Показать меню",
		Coordinates:   "📍 Поделится с кординатами местоположения",
		Radius:        "🔍 Настроить радиус поиска",
	}
}

type PlaceCategories struct {
	CafeAndRestaurants string

	CoffeeAndTea string

	JapaneseFood string

	ConfectioneryStores string

	GymAndStudios string

	SightsAndNature string

	ChurchesAndMosques string

	HealthAndBeautyServices string

	CityParks string

	Hotels string

	Beaches string

	Shops string

	Bars string
}

func NewPlaceCategories() *PlaceCategories {
	return &PlaceCategories{
		CafeAndRestaurants:      "🍽️ Кафе-Рестораны",
		CoffeeAndTea:            "🍵 Кофейная-Чайная",
		JapaneseFood:            "🍣 Японская кухня",
		ConfectioneryStores:     "🍰 Кондитерские магазины",
		GymAndStudios:           "🏋️‍♀️ Тренажерный зал и студии",
		SightsAndNature:         "🏛️ Достопремечательности и природа",
		ChurchesAndMosques:      "⛪️ Церквки-Мечети",
		HealthAndBeautyServices: "💆‍♀️ Услуги для здоровья и красоты",
		CityParks:               "🌳 Городские парки",
		Hotels:                  "🏨 Отели",
		Beaches:                 "🏖️ Пляжи",
		Shops:                   "🛍️ Магазины",
		Bars:                    "🍻 Бары",
	}
}

type PlaceCategoriesCode struct {
	CafeAndRestaurantsCode string

	CoffeeAndTeaCode string

	JapaneseFoodCode string

	ConfectioneryStoresCode string

	GymAndStudiosCode string

	SightsAndNatureCode string

	ChurchesAndMosquesCode string

	HealthAndBeautyServicesCode string

	CityParksCode string

	HotelsCode string

	BeachesCode string

	ShopsCode string

	BarsCode string
}

func NewPlaceCategoriesCode() *PlaceCategoriesCode {
	return &PlaceCategoriesCode{
		CafeAndRestaurantsCode:      "13065",
		CoffeeAndTeaCode:            "13032",
		JapaneseFoodCode:            "13263",
		ConfectioneryStoresCode:     "17060",
		GymAndStudiosCode:           "18077",
		SightsAndNatureCode:         "16020",
		ChurchesAndMosquesCode:      "12101",
		HealthAndBeautyServicesCode: "11061",
		CityParksCode:               "16032",
		HotelsCode:                  "19014",
		BeachesCode:                 "16003",
		ShopsCode:                   "17043",
		BarsCode:                    "13003",
	}
}
