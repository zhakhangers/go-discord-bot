package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Store Bot API Tokens:
var (
	OpenWeatherToken string
	BotToken         string
)

func Run() {
	// Create new Discord Session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	// Cleanly close down the Discord session.
	defer discord.Close()

	// Add event handler for general messages
	discord.AddHandler(NewBotHandler)

	// Run until code is terminated
	fmt.Println("Bot running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

}
