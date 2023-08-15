package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	utils "github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/cache"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Response struct {
	Results []Location `json:"results"`
}

type Location struct {
	Name     string `json:"name"`
	Distance int    `json:"distance"`
	FsqID    string `json:"fsq_id"`
	Location struct {
		Address     string `json:"address"`
		CrossStreet string `json:"cross_street"`
		Locality    string `json:"locality"`
		Region      string `json:"region"`
		Country     string `json:"country"`
		PostCode    string `json:"postcode"`
	} `json:"location"`
	Geocodes struct {
		Main struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"main"`
	} `json:"geocodes"`
}

func GetNearbyPlaces(limitSearch, limitPhotos int64, categoryID string, bot *tgbotapi.BotAPI,
	update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	latitude, longitude := GetCoordinates(update.Message.Chat.ID)
	queryURL := buildQueryURL(limitSearch, latitude, longitude, categoryID, update.Message.Chat.ID)
	locations, err := searchPlaces(queryURL)
	if err != nil {
		utils.Fatal(err)
	}

	var messageText bytes.Buffer
	for _, location := range locations {
		messageText.WriteString(fmt.Sprintf("Название: %s\n", location.Name))
		messageText.WriteString(fmt.Sprintf("Адрес: %s\n", location.Location.Address))
		messageText.WriteString(fmt.Sprintf("Расстояние: %d \n", location.Distance))
		photoURL, err := getPhotosPlaces(limitPhotos, location.FsqID)
		if err != nil {
			utils.Error(err.Error())
		} else {
			messageText.WriteString(fmt.Sprintf("Адрес Фота: %s\n", photoURL))
		}
		messageText.WriteString("\n")
	}

	result := messageText.String()

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подробнее о местe", "button1"),
			tgbotapi.NewInlineKeyboardButtonData("Показать на карте", "button2"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть отзывы места", "button3"),
			tgbotapi.NewInlineKeyboardButtonData("Оставить отзыв о месте", "button4"),
		),
	)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сохранить место"),
			tgbotapi.NewKeyboardButton("Добавить в избранное"),
		),

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы также можете сохранить места или добавить их в список избранных мест.")
	msg2.ReplyMarkup = replyKeyboard
	bot.Send(msg2)

	errCh := make(chan string)

	for upd := range updates {
		if upd.CallbackQuery != nil && upd.CallbackQuery.Data != "" {
			switch upd.CallbackQuery.Data {
			case "button1":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места:")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)
				waitInputUser(place, errCh, updates)
				detailAboutPlace(bot, update.Message.Chat.ID, <-place, locations)
				break
			case "button2":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места из списка предыдущих мест:")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)
				waitInputUser(place, errCh, updates)
				showInMap(<-place, bot, update.Message.Chat.ID, locations)

			case "button3":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места для получения отзывов о нём: ")
				bot.Send(msg)

				place := make(chan string)
				waitInputUser(place, errCh, updates)
				checkReviewOfThePlace(bot, update.Message.Chat.ID, place, locations)

			case "button4":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите место в котором хотите оставить свой отзыв: ")
				bot.Send(msg)

				place := make(chan string)
				waitInputUser(place, errCh, updates)
				leftReviewOfThePlace(bot, update.Message.Chat.ID, updates, update, place, locations)
			}

		} else if update.Message != nil && update.Message.Text != "" {
			switch upd.Message.Text {
			case "Сохранить место":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места которое хотите сохранить.")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)
				waitInputUser(place, errCh, updates)
				savePlace(bot, update.Message.Chat.ID, <-place, locations)

			case "Добавить в избранное":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места которое хотите добавить в израбнное.")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)
				waitInputUser(place, errCh, updates)
				saveFavoritePlace(bot, update.Message.Chat.ID, <-place, locations)

			case "Назад":
				msg.Text = "Выберите действие: "
				msg.ReplyMarkup = createNeedAction()
				bot.Send(msg)
				return
			}
		}
	}
}

func waitInputUser(place chan string, errCh chan string, updates tgbotapi.UpdatesChannel) {
	go func() {
		for {
			select {
			case input := <-updates:
				if input.Message != nil {
					place <- input.Message.Text
					return
				}
			case <-time.After(time.Second * 30):
				place <- ""
				return
			case err := <-errCh:
				if err == "Произошла ошибка" {
					place <- ""
					return
				}
			}
		}
	}()
}

func detailAboutPlace(bot *tgbotapi.BotAPI, chatID int64, placeName string, locations []Location) {
	var messageText bytes.Buffer
	found := false
	for _, location := range locations {
		if placeName == location.Name {
			messageText.WriteString(fmt.Sprintf("Местонохождение: %s\n", location.Location.Locality))
			messageText.WriteString(fmt.Sprintf("Регион: %s\n", location.Location.Region))
			messageText.WriteString(fmt.Sprintf("Перекресток улицы: %s\n", location.Location.CrossStreet))
			messageText.WriteString(fmt.Sprintf("Страна: %s\n", location.Location.Country))
			messageText.WriteString(fmt.Sprintf("Почтовый индекс: %s\n", location.Location.PostCode))
			messageText.WriteString("\n")
			found = true
		}
	}

	if !found {
		msg := tgbotapi.NewMessage(chatID, "Такого места нет в списке,будьте внимательны.")
		bot.Send(msg)
	}

	result := messageText.String()
	msg := tgbotapi.NewMessage(chatID, result)
	if _, err := bot.Send(msg); err != nil {
		return
	}
}

func buildQueryURL(limit int64, latitude, longitude float64, categoryID string, userID int64) string {
	dbRadius := database.NewRadiusSearchImpl(database.DB)
	err, radius := dbRadius.GetRadiusSearch(context.Background(), userID)
	if err != nil {
		radius = 100000
		utils.Error("Failed to get radius user's")
	}

	queryParams := url.Values{}
	queryParams.Set("ll", fmt.Sprintf("%f,%f", latitude, longitude))
	queryParams.Set("categories", categoryID)
	queryParams.Set("client_id", os.Getenv("CLIENT_ID"))
	queryParams.Set("radius", strconv.FormatFloat(radius, 'f', 0, 64))
	queryParams.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	queryParams.Set("oauth_token", os.Getenv("API_TOKEN"))
	queryParams.Set("limit", strconv.FormatInt(limit, 10))
	queryParams.Set("sort", "distance")
	queryParams.Set("open_now", "true")

	return fmt.Sprintf("%s?%s", os.Getenv("API_URLSEARCHPLACE"), queryParams.Encode())
}

func searchPlaces(queryURL string) ([]Location, error) {
	c := cache.NewCache()
	if cacheData, found := c.Get(queryURL); found {
		return cacheData.([]Location), nil
	}

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		utils.Error("Failed to send request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", os.Getenv("API_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Error("Failed to get response: %v", err)
	}

	defer res.Body.Close()

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	c.Set(queryURL, response.Results, time.Minute*5)

	return response.Results, nil
}
