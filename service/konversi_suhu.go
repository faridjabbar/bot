package service

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func convertTemperature(temperature float64, fromUnit, toUnit string) float64 {
	var convertedTemperature float64

	if fromUnit == toUnit {
		return temperature
	}

	if fromUnit == "Celsius" {
		if toUnit == "Fahrenheit" {
			convertedTemperature = (temperature * 9 / 5) + 32
		} else if toUnit == "Kelvin" {
			convertedTemperature = temperature + 273.15
		}
	} else if fromUnit == "Fahrenheit" {
		if toUnit == "Celsius" {
			convertedTemperature = (temperature - 32) * 5 / 9
		} else if toUnit == "Kelvin" {
			convertedTemperature = (temperature + 459.67) * 5 / 9
		}
	} else if fromUnit == "Kelvin" {
		if toUnit == "Celsius" {
			convertedTemperature = temperature - 273.15
		} else if toUnit == "Fahrenheit" {
			convertedTemperature = (temperature * 9 / 5) - 459.67
		}
	}

	return convertedTemperature
}

func HandleTemperatureConversion(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	commandParts := strings.Fields(update.Message.CommandArguments())
	if len(commandParts) != 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gunakan perintah dengan format yang benar, misalnya: /converttemperature 30 Celsius Fahrenheit")
		bot.Send(msg)
		return
	}

	temperature, err := strconv.ParseFloat(commandParts[0], 64)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nilai suhu tidak valid. Gunakan angka.")
		bot.Send(msg)
		return
	}

	fromUnit := commandParts[1]
	toUnit := commandParts[2]

	convertedTemperature := convertTemperature(temperature, fromUnit, toUnit)

	msgText := fmt.Sprintf("%.2f %s sama dengan %.2f %s.", temperature, fromUnit, convertedTemperature, toUnit)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}
