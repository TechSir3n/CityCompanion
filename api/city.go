package api

import (
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/TechSir3n/CityCompanion/database"
)

type Context struct {
	Context struct {
		GeoBounds struct {
			Circle struct {
				Center struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"center"`
			} `json:"circle"`
		} `json:"geo_bounds"`
	} `json:"context"`
}

func GetCordinatesByCity(cityName string) (float64, float64, error) {
	if cityName == "" {
		return 0.0, 0.0, nil
	}

	url := "https://api.foursquare.com/v3/places/search?near=" + cityName

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", os.Getenv("API_TOKEN"))

	if err != nil {
		return 0.0, 0.0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0.0, 0.0, err
	}

	defer res.Body.Close()

	var context Context
	err = json.NewDecoder(res.Body).Decode(&context)
	if err != nil {
		return 0.0, 0.0, err
	}

	latitude := context.Context.GeoBounds.Circle.Center.Latitude
	longitude := context.Context.GeoBounds.Circle.Center.Longitude
	return latitude, longitude, nil
}
