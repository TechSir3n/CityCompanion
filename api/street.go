package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TechSir3n/CityCompanion/assistance"
	"github.com/TechSir3n/CityCompanion/database"
)

type GeocodingResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted"`
	} `json:"results"`
}

func getUserStreet(userID int64) (string, error) {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background(), userID)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%f,%f&key=%s",
		latitude, longitude, os.Getenv("API_OPENCAGEDATA"))

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var geocodingResponse GeocodingResponse
	err = json.Unmarshal(body, &geocodingResponse)
	if err != nil {
		return "", err
	}

	var address string
	if len(geocodingResponse.Results) > 0 {
		address = geocodingResponse.Results[0].FormattedAddress
	} else {
		return "", err
	}

	return address, nil
}

func getCoordinatesByStreet(streetName string) (float64, float64, error) {
	var street string
	if assistance.IsRussianWord(streetName) {
		street, _ = translateRussianToEnglish(streetName)
	} else {
		street = streetName
	}

	substr := strings.Split(street, ",")
	queryParams := url.Values{}
	queryParams.Set("near", fmt.Sprintf("%s,%s", substr[0], substr[len(substr)-1]))

	apiURL := fmt.Sprintf("https://api.foursquare.com/v3/places/search?%s", queryParams.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0.0, 0.0, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", os.Getenv("API_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0.0, 0.0, err
	}

	defer res.Body.Close()

	var cordinates Context
	err = json.NewDecoder(res.Body).Decode(&cordinates)
	if err != nil {
		return 0.0, 0.0, err
	}

	latitude := cordinates.Context.GeoBounds.Circle.Center.Latitude
	longitude := cordinates.Context.GeoBounds.Circle.Center.Longitude

	return longitude, latitude, nil
}

func GetCoordinates(userID int64) (float64, float64) {
	userLocation := database.NewUserLocationImpl(database.DB)
	err, latitude, longitude := userLocation.GetUserLocation(context.Background(), userID)
	if err != nil {
		return 0.0, 0.0
	}

	return latitude, longitude
}
