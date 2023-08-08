package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/TechSir3n/CityCompanion/database"
)

type Photo struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func getPhotosPlaces(limit int64, fsqID string) ([]string, error) {
	queryURL := fmt.Sprintf("https://api.foursquare.com/v3/places/%s/photos?limit=%d", fsqID, limit)
	req, err := http.NewRequest("GET", queryURL, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("classifications", "outdoor")
	req.Header.Add("sort", "POPULAR")
	req.Header.Add("Authorization", os.Getenv("API_TOKEN"))

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var photos []Photo
	err = json.NewDecoder(res.Body).Decode(&photos)
	if err != nil {
		return nil, err
	}

	var photoURLs = make([]string,len(photos))
	for _, photo := range photos {
		photoURL := fmt.Sprintf("%s%s%s", photo.Prefix, "original", photo.Suffix)
		photoURLs = append(photoURLs,photoURL)
	}

	return photoURLs, nil
}
