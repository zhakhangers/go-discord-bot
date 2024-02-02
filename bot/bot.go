package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Storing Bot API Tokens:
var (
	OpenWeatherToken     string
	BotToken             string
	GoogleTranslateToken string
)

// Run() initializes and starts the bot.
func Run() {
	// Create a new Discord session.
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	// Opens a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	// Cleanly closes down the Discord session when finished.
	defer discord.Close()

	// Adds an event handler for general messages.
	discord.AddHandler(NewBotHandler)

	// Runs the bot until it is terminated.
	fmt.Println("Bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

}
