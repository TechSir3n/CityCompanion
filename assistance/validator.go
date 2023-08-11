package assistance

import (
	"strconv"
	"strings"
	"unicode"
)

func IsRussianWord(word string) bool {
	for _, c := range word {
		if !unicode.Is(unicode.Cyrillic, c) {
			return false
		}
	}
	return true
}

func IsRadiusCorrect(radius string) bool {
	if strings.HasSuffix(radius, "км") {
		radius = strings.TrimSuffix(radius, "км")
		num, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			return false
		} else if num < 0 {
			return false
		}
		return true
	} else if strings.HasSuffix(radius, "м") {
		radius = strings.TrimSuffix(radius, "м")
		num, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			return false
		} else if num < 0 {
			return false
		}
		return true
	}
	return false
}

func IsRightFormat(street string) bool {
	if !strings.Contains(street, ",") {
		return false
	}

	street = strings.TrimSpace(street)
	parts := strings.Split(street, ",")
	if len(parts) != 2 {
		return false
	}

	if parts[0] == "" || parts[1] == "" {
		return false
	}

	for _, ch := range street {
		if !unicode.IsLetter(ch) && !unicode.IsSpace(ch) && ch != ',' {
			return false
		}
	}

	return true
}
