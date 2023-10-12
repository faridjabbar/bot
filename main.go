package main

import (
	c "bot/mod/configuration"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	configuration, err := c.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	bot, err := tgbotapi.NewBotAPI(""+configuration.ApiKey+"")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Bot %s berhasil diinisialisasi", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Halo!")
				bot.Send(msg)
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ini adalah bot bantuan. Anda dapat melakukan perintah-perintah berikut:\n/start - Memulai bot\n/help - Menampilkan bantuan")
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Perintah tidak dikenal. Gunakan /help untuk melihat perintah yang tersedia.")
			bot.Send(msg)
		}
	}
}
