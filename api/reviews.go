package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/TechSir3n/CityCompanion/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"time"
)

func leftReviewOfThePlace(bot *tgbotapi.BotAPI, chatID int64, updates tgbotapi.UpdatesChannel,
	update tgbotapi.Update, place chan string, locations []Location) {
	db := database.NewReviewPlacesImp(database.DB)
	found := false
	locName := <-place
	errCh:=make(chan string)
	for _, location := range locations {
		if locName == location.Name {
			review := tgbotapi.NewMessage(chatID, "Отлично,теперь можете оставить свой отзыв")
			bot.Send(review)

			ratingMsg := tgbotapi.NewMessage(chatID, "Оцените место по шкале от 1 до 5 : ")
			bot.Send(ratingMsg)

			inputRating := make(chan string)
			waitInputUser(inputRating,errCh, updates)
			rating := <-inputRating

			if rating != "" {
				nameInput := tgbotapi.NewMessage(chatID, "Ввведите ваше имя: ")
				bot.Send(nameInput)
			}

			inputName := make(chan string)
			waitInputUser(inputName,errCh, updates)

			name := <-inputName
			if name != "" {
				commentInput := tgbotapi.NewMessage(chatID, "Теперь,введите ваш комментарий: ")
				bot.Send(commentInput)
			}

			inputComment := make(chan string)
			waitInputUser(inputComment,errCh, updates)
			comment := <-inputComment

			if comment != "" {
				msg := tgbotapi.NewMessage(chatID, "Готово !")
				bot.Send(msg)
			}

			found = true
			r, _ := strconv.Atoi(rating)
			if err := db.SaveReview(context.Background(), location.Name, location.Location.Address,
				comment, name, r); err != nil {
				return
			} else {
				msg := tgbotapi.NewMessage(chatID, "Отзыв успешно получен,благодарим вас за него.")
				bot.Send(msg)
			}
		}
	}


	if !found {
		msg := tgbotapi.NewMessage(chatID, "Вы ввели место, которого нет в списке.")
		bot.Send(msg)
	}
}

func checkReviewOfThePlace(bot *tgbotapi.BotAPI, chatID int64, place chan string, locations []Location) {
	locName := <-place
	found := false
	db := database.NewReviewPlacesImp(database.DB)
	var messageText bytes.Buffer
	for _, location := range locations {
		messageText.Reset()
		if location.Name == locName {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if comments, usernames, ratings, created, err := db.GetReview(ctx, locName); err != nil {
				msg := tgbotapi.NewMessage(chatID, "Простите,отзывов об этом заведение еще нет, или их не удалось получить.")
				bot.Send(msg)
			} else {
				for i := 0; i < len(usernames); i++ {
					messageText.WriteString(fmt.Sprintf("Имя: %s\n", usernames[i]))
					messageText.WriteString(fmt.Sprintf("Комментарий: %s\n", comments[i]))
					messageText.WriteString(fmt.Sprintf("Рейнтинг: %d\n", ratings[i]))
					input := created[i]
					t, err := time.Parse(time.RFC3339Nano, input)
					if err != nil {
						fmt.Println("Error parsing time:", err)
						return
					}
					output := t.Format("2006-01-02 15:04")
					messageText.WriteString(fmt.Sprintf("Дата: %s\n", output))
					messageText.WriteString("\n")
				}
			}

			result := messageText.String()
			msg := tgbotapi.NewMessage(chatID, result)
			if _, err := bot.Send(msg); err != nil {
				fmt.Printf("Error send:  %v", err)
			}

			found = true
		}
	}

	if !found {
		msg := tgbotapi.NewMessage(chatID, "Вы ввели место, которого нет в списке.")
		bot.Send(msg)
	}
}
