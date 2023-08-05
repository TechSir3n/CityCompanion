package api

import (
	"context"
	"encoding/json"
	"fmt"
	utils "github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

type Response struct {
	Results []Location `json:"results"`
}

type Location struct {
	Name     string `json:"name"`
	Distance int    `json:"distance"`
	FsqID    string `json:"fsq_id"`
	Location struct {
		Address string `json:"address"`
	} `json:"location"`
}

func GetNearbyPlaces(limitSearch, limitPhotos int64, categoryID string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	latitude, longitude := GetCoordinates()
	queryURL := buildQueryURL(limitSearch, latitude, longitude, categoryID)
	locations, err := searchPlaces(queryURL)
	if err != nil {
		utils.Fatal(err)
	}

	var messageText string
	for _, location := range locations {
		messageText += fmt.Sprintf("Название: %s\n", location.Name)
		messageText += fmt.Sprintf("Адрес: %s\n", location.Location.Address)
		messageText += fmt.Sprintf("Расстояние: %d м\n", location.Distance)
		photoURL, err := getPhotosPlaces(limitPhotos, location.FsqID)
		if err != nil {
			utils.Error("Here:", err.Error())
		} else {
			messageText += fmt.Sprintf("Адресс Фота: %s м\n", photoURL)
		}
		messageText += "\n"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	bot.Send(msg)
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

	sort.Slice(response.Results, func(i, j int) bool {
		return response.Results[i].Distance < response.Results[j].Distance
	})

	return response.Results, nil
}
