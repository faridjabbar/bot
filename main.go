package main

import (
	c "bot-telegram/configuration"
	"bot-telegram/service"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	configuration, err := c.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	bot, err := tgbotapi.NewBotAPI("" + configuration.ApiKey + "")
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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ini adalah bot bantuan. Anda dapat melakukan perintah-perintah berikut:\n/start - untuk memunculkan kata Halo!\n/help - Menampilkan bantuan\n/converttosticker - merubah gambar jadi stiker (beta version v:)"+
					"\n/converttime - merubah waktu\n/converttemperature - Converted Temperature\n/kritiksaran untuk memberikan kritik dan saran")
				bot.Send(msg)
			case "converttosticker":
				service.HandleConvertToSticker(&update, bot, updates, configuration)
			case "converttime":
				convertedValue := service.ConvertTime(3600, "detik", "jam")
				msgText := fmt.Sprintf("%d %s sama dengan %s %s.", 3600, "detik", convertedValue, "jam")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				bot.Send(msg)
			case "converttemperature":
				service.HandleTemperatureConversion(&update, bot)
			case "kritiksaran":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Silakan masukkan kritik dan saran Anda:")
				bot.Send(msg)

				nextUpdate := <-updates
				if nextUpdate.Message == nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Maaf, pesan tidak valid. Perintah dibatalkan.")
					bot.Send(msg)
				} else {
					kritikDanSaran := nextUpdate.Message.Text
					kritikDanSaranMessage := tgbotapi.NewMessage(2050877699, "Kritik dan Saran:\n"+kritikDanSaran)
					bot.Send(kritikDanSaranMessage)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terima kasih atas kritik dan sarannya. Pesan Anda telah terkirim.")
					bot.Send(msg)
				}

			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Perintah tidak dikenal. Gunakan /help untuk melihat perintah yang tersedia.")
			bot.Send(msg)
		}
	}
}
