package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TechSir3n/CityCompanion/database"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


type GeocodingResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted"`
	} `json:"results"`
}

func GetUserStreet() string {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background())
	if err != nil {
		log.Fatalf("Failed to get user gelocation: %v", err)
	}

	url := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%f,%f&key=%s",
		latitude, longitude, os.Getenv("API_OPENCAGEDATA"))

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return ""
	}

	var geocodingResponse GeocodingResponse
	err = json.Unmarshal(body, &geocodingResponse)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return ""
	}

	var address string
	if len(geocodingResponse.Results) > 0 {
		address = geocodingResponse.Results[0].FormattedAddress
		fmt.Println(address)
	} else {
		fmt.Println("Адрес не найден")
	}

	return address
}

func GetCoordinates() (float64, float64) {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background())
	if err != nil {
	
	}

	return latitude, longitude
}
