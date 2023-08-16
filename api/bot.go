package api

import (
	"fmt"
	"log"
	"os"

	_ "github.com/TechSir3n/CityCompanion/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func TGBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error[CreateButtonFunc] : %v", err)
	}

	bot.Debug = false
	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err!=nil { 
		log.Fatalf("Error[GetUpdatesChan]: %v",err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		commandsBot(bot, update, updates)
		categoriesPlace(bot, update, updates)
	}
}
