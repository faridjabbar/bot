package service

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ConvertTime(timeValue int, fromUnit, toUnit string) string {
	var convertedValue float64

	switch fromUnit {
	case "detik":
		switch toUnit {
		case "menit":
			convertedValue = float64(timeValue) / 60
		case "jam":
			convertedValue = float64(timeValue) / 3600
		case "hari":
			convertedValue = float64(timeValue) / 86400
		case "minggu":
			convertedValue = float64(timeValue) / 604800
		case "tahun":
			convertedValue = float64(timeValue) / 31536000
		}
	case "menit":
		minutes := float64(timeValue)
		switch toUnit {
		case "detik":
			convertedValue = minutes * 60
		case "jam":
			convertedValue = minutes / 60
		case "hari":
			convertedValue = minutes / 1440
		case "minggu":
			convertedValue = minutes / 10080
		case "tahun":
			convertedValue = minutes / 525600
		}
	case "jam":
		hours := float64(timeValue)
		switch toUnit {
		case "detik":
			convertedValue = hours * 3600
		case "menit":
			convertedValue = hours * 60
		case "hari":
			convertedValue = hours / 24
		case "minggu":
			convertedValue = hours / 168
		case "tahun":
			convertedValue = hours / 8760
		}
	}

	return strconv.FormatFloat(convertedValue, 'f', -1, 64)
}

func HandleTimeConversion(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	commandParts := strings.Fields(update.Message.CommandArguments())
	if len(commandParts) != 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gunakan perintah dengan format yang benar, misalnya: /converttime 3600 detik jam")
		bot.Send(msg)
		return
	}

	timeValue, err := strconv.Atoi(commandParts[0])
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nilai waktu tidak valid. Gunakan angka.")
		bot.Send(msg)
		return
	}

	fromUnit := commandParts[1]
	toUnit := commandParts[2]

	convertedValue := ConvertTime(timeValue, fromUnit, toUnit)

	msgText := fmt.Sprintf("%s %s sama dengan %s %s.", strconv.Itoa(timeValue), fromUnit, convertedValue, toUnit)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}
