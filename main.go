package main

import (
	"log"
	"os"
	"test-project-bot/bot"
)

func main() {
	// Load environment variables
	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Must set Discord token as env variable: BOT_TOKEN")
	}
	openWeatherToken, ok := os.LookupEnv("OPENWEATHER_TOKEN")
	if !ok {
		log.Fatal("Must set Open Weather token as env variable: OPENWEATHER_TOKEN")
	}
	googleTranslateToken, ok := os.LookupEnv(("GOOGLETRANSLATE_TOKEN"))
	if !ok {
		log.Fatal("Must set Open Weather token as env variable: GOOGLETRANSLATE_TOKEN")
	}

	// Save API keys & start bot
	bot.BotToken = botToken
	bot.OpenWeatherToken = openWeatherToken
	bot.GoogleTranslateToken = googleTranslateToken
	bot.Run()
}
