package service

import (
	"io"
	"log"
	"net/http"

	c "bot-telegram/configuration"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleConvertToSticker(update *tgbotapi.Update, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, configuration c.Configuration) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Silakan kirim gambar yang ingin Anda konversi menjadi stiker.")
	bot.Send(msg)

	nextUpdate := <-updates
	if nextUpdate.Message == nil || nextUpdate.Message.Photo == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Maaf, tidak ada gambar yang dikirim. Perintah dibatalkan.")
		bot.Send(msg)
	} else {
		photo := (*nextUpdate.Message.Photo)[0]

		fileConfig := tgbotapi.FileConfig{
			FileID: photo.FileID,
		}
		file, err := bot.GetFile(fileConfig)
		if err != nil {
			log.Println("Gagal mendapatkan file:", err)
			return
		}

		fileURL := "https://api.telegram.org/file/bot" + configuration.ApiKey + "/" + file.FilePath
		response, err := http.Get(fileURL)
		if err != nil {
			log.Println("Gagal mengunduh gambar:", err)
			return
		}
		defer response.Body.Close()

		imageData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Gagal membaca gambar:", err)
			return
		}

		sticker := tgbotapi.NewStickerUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
			Name:  "sticker.png",
			Bytes: imageData,
		})

		bot.Send(sticker)
	}

}
