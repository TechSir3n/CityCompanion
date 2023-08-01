package api

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleQuery(update tgbotapi.Update) {
	fmt.Println("We here")

	callback:=tgbotapi.NewCallback(update.CallbackQuery.ID,update.CallbackQuery.Data)
	fmt.Println("Text:",callback.Text)

	switch update.CallbackQuery.Data {
	case "action1":
		fmt.Println("I here ")

	case "action2":

		break

	case "action3":

		break

	case "action4":

		break

	case "action5":

		break

	case "action6":

		break

	case "action7":

		break

	case "action8":

		break

	case "action9":

		break

	case "action10":

		break

	case "action11":

		break

	default:
		break
	}
}
