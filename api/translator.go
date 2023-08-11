package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Translation struct {
	Text string `json:"text"`
}

type Word struct {
	Translations []Translation `json:"translations"`
}

func translateRussianToEnglish(word string) (string, error) {
	queryStr := url.QueryEscape(word)
	URL := fmt.Sprintf("https://linguee-api.fly.dev/api/v2/translations?query=%s&src=ru&dst=en&guess_direction=true&follow_corrections=always", queryStr)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var translate []Word

	err = json.NewDecoder(res.Body).Decode(&translate)
	if err != nil {
		return "", err
	}

	var englishStr string
	for _, t := range translate {
		englishStr = t.Translations[0].Text
	}

	return englishStr, nil
}
