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
	latitude, longitude := GetCoordinates()
	queryURL := buildQueryURL(limitSearch, latitude, longitude, categoryID)
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
			tgbotapi.NewInlineKeyboardButtonData("Подробнее о местах", "button1"),
			tgbotapi.NewInlineKeyboardButtonData("Показать на карте", "button2"),
		),
	)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сохранить место"),
			tgbotapi.NewKeyboardButton("Добавить в избранное"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы также можете сохранить места или добавить их в список избранных мест.")
	msg2.ReplyMarkup = replyKeyboard
	bot.Send(msg2)

	var detailInfo string
	for upd := range updates {
		if upd.CallbackQuery != nil && upd.CallbackQuery.Data != "" {
			switch upd.CallbackQuery.Data {
			case "button1":
				for _, location := range locations {
					detailInfo += fmt.Sprintf("Перекресток: %s\n", location.Location.CrossStreet)
					detailInfo += fmt.Sprintf("Местонохождение: %s\n", location.Location.Locality)
					detailInfo += fmt.Sprintf("Регион: %s\n", location.Location.Region)
					detailInfo += fmt.Sprintf("Страна: %s\n", location.Location.Country)
					detailInfo += fmt.Sprintf("Почтовый индекс: %s\n", location.Location.PostCode)
					detailInfo += "\n"
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, detailInfo)
				if _, err := bot.Send(msg); err != nil {
					return
				}
				break
			case "button2":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места из списка предыдущих мест:")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)

				go func() {
					for {
						select {
						case input := <-updates:
							if input.Message != nil {
								place <- input.Message.Text
								break
							}
						case <-time.After(30 * time.Second):
							break
						}
					}
				}()

				found := false
				enteredName := <-place

				for _, location := range locations {
					if enteredName == location.Name {
						msg := tgbotapi.NewLocation(update.Message.Chat.ID,
							location.Geocodes.Main.Latitude, location.Geocodes.Main.Longitude)
						if _, err := bot.Send(msg); err != nil {
							return
						}
						found = true
						return
					}
				}

				if !found {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Простите,но такого мест нет в спискe.")
					if _, err := bot.Send(msg); err != nil {
						return
					}
				}
			}
		} else if update.Message != nil && update.Message.Text != "" {
			switch upd.Message.Text {
			case "Сохранить место":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места которое хотите сохранить.")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)
				go func() {
					for {
						select {
						case input := <-updates:
							if input.Message != nil {
								place <- input.Message.Text
								break
							}
						case <-time.After(time.Second * 30):
							break
						}
					}
				}()

				placeName := <-place
				savePlace(bot, update.Message.Chat.ID, placeName, locations)
			case "Добавить в избранное":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название места которое хотите добавить в израбнное.")
				if _, err := bot.Send(msg); err != nil {
					return
				}

				place := make(chan string)
				go func() {
					for {
						select {
						case input := <-updates:
							if input.Message != nil {
								place <- input.Message.Text
								break
							}
						case <-time.After(time.Second * 30):
							break
						}
					}
				}()
				placeName := <-place
				saveFavoritePlace(bot, update.Message.Chat.ID, placeName, locations)
			}
		}
	}
}

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

func buildQueryURL(limit int64, latitude, longitude float64, categoryID string) string {
	dbRadius := database.NewRadiusSearchImpl(database.DB)
	err, radius := dbRadius.GetRadiusSearch(context.Background())
	if err != nil {
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

	return fmt.Sprintf("%s?%s", os.Getenv("API_URLSEARCHPLACE"), queryParams.Encode())
}

func searchPlaces(queryURL string) ([]Location, error) {
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
	
	return response.Results, nil
}
