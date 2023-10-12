package main

import (
	c "bot/mod/configuration"
	"io/ioutil"
	"log"
	"net/http"

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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ini adalah bot bantuan. Anda dapat melakukan perintah-perintah berikut:\n/start - Memulai bot\n/help - Menampilkan bantuan\n/converttosticker - merubah gambar jadi stiker")
				bot.Send(msg)
			case "converttosticker":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Silakan kirim gambar yang ingin Anda konversi menjadi stiker.")
				bot.Send(msg)

				// Membaca pesan berikutnya yang harus berupa gambar
				nextUpdate := <-updates
				if nextUpdate.Message == nil || nextUpdate.Message.Photo == nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Maaf, tidak ada gambar yang dikirim. Perintah dibatalkan.")
					bot.Send(msg)
				} else {
					// Mengambil foto pertama dari pesan
					photo := (*nextUpdate.Message.Photo)[0]

					// Mengunduh foto ke server Anda
					fileConfig := tgbotapi.FileConfig{
						FileID: photo.FileID,
					}
					file, _ := bot.GetFile(fileConfig)

					// Mengunduh gambar ke server Anda
					fileURL := "https://api.telegram.org/file/bot" + configuration.ApiKey + "/" + file.FilePath
					response, err := http.Get(fileURL)
					if err != nil {
						log.Println("Gagal mengunduh gambar:", err)
						return
					}
					defer response.Body.Close()

					// Membaca gambar sebagai byte slice
					imageData, err := ioutil.ReadAll(response.Body)
					if err != nil {
						log.Println("Gagal membaca gambar:", err)
						return
					}

					// Mengirim gambar sebagai stiker dengan ukuran normal
					sticker := tgbotapi.NewStickerUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
						Name:  "sticker.png",
						Bytes: imageData,
					})

					bot.Send(sticker)
				}
			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Perintah tidak dikenal. Gunakan /help untuk melihat perintah yang tersedia.")
			bot.Send(msg)
		}
	}
}
