package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
	routes map[*Route]interface{} // A map to store routes and their associated functions.
}

type Route struct {
	name        string // Name of the route.
	description string // Description of the route.
}

type contextKey string

const (
	cmdKey    contextKey = "cmd" // Key for storing the command in the context.
	prefixKey contextKey = "!"   // Key for storing the command prefix in the context.
)

// handleRoute handles incoming commands and executes corresponding routes.
func (r *Router) handleRoute(c *Client) {
	cmd := c.ctx.Value("cmd").(string) // Retrieves the command from the context.

	if strings.HasPrefix(cmd, "help") {
		r.help(c) // Handles the "help" command separately.
	}

	// Iterates through registered routes to find a matching one.
	for routeName, routeFunc := range r.routes {
		if strings.HasPrefix(cmd, routeName.name) {
			log.Printf("Route found based on command | %s -> %s", cmd, routeName)
			routeFunc.(func())() // Execute the associated route function.
		}
	}
}

func (r *Router) help(c *Client) {
	var res string

	// Iterate through registered routes and collect their descriptions.
	for route := range r.routes {
		res = res + "\n" + route.description
	}

	// Send the collected route descriptions as a response.
	if msg, err := c.session.ChannelMessageSend(c.message.ChannelID, res); err != nil {
		log.Printf("error sending message | %v | %v", msg, err)
	}
}

func NewBotHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := "!"
	ctx := context.WithValue(context.Background(), prefixKey, prefix)
	client := Client{ctx, s, m}

	cmd := client.message.Content
	client.ctx = context.WithValue(client.ctx, cmdKey, cmd)

	// Check if the message author is the bot itself.
	if client.message.Author.ID == client.session.State.User.ID {
		return
	}

	// Check if the message has the expected prefix.
	if !strings.HasPrefix(client.message.Content, prefix) {
		return
	}

	log.Print("Prefix hit")
	cmd = strings.TrimPrefix(client.message.Content, prefix)
	client.ctx = context.WithValue(client.ctx, "cmd", cmd)

	// Define routes and create a router instance.
	routes := map[*Route]interface{}{
		&Route{"weather", "`!weather help` - Weather info from Open Weather Maps "}: client.Weather,
		&Route{"translate", "`!translate help` - translate "}:                       client.Translate,
	}

	router := &Router{routes: routes}

	// Handle the incoming command using the router.
	router.handleRoute(&client)

}

func (client *Client) Weather() {
	log.Printf("go.bot.Weather request recieved")

	cmd := client.ctx.Value("cmd").(string)
	cmd = strings.TrimPrefix(cmd, "weather ")
	client.ctx = context.WithValue(client.ctx, "cmd", cmd)

	// Define weather-related routes and create a router instance.
	routes := map[*Route]interface{}{
		&Route{"get", "`!weather get [location]` - gets the weather for a location"}: client.GetWeather,
	}

	weatherRouter := &Router{routes: routes}

	// Handle the weather-related command using the weather router.
	weatherRouter.handleRoute(client)
}

func (client *Client) GetWeather() {
	log.Printf("go.bot.handler.Weather.Get request recieved")

	// Extracts the location from the command.
	cmd := client.ctx.Value("cmd").(string)
	cmd = strings.TrimPrefix(cmd, "get ")
	client.ctx = context.WithValue(client.ctx, "cmd", cmd)

	// Performs the weather request asynchronously.
	go client.processWeatherRequest(cmd)
}

func (client *Client) processWeatherRequest(cmd string) {
	// Creates a weather client and retrieve weather data.
	weatherClient := NewWeatherClient(client.ctx)
	res, err := weatherClient.GetWeatherByLocation(cmd)
	if err != nil {
		log.Printf("error getting weather: %s | %v", cmd, err)
		embeddedMsg := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{},
			Title:       "There is no such location",
			Timestamp:   time.Now().Format(time.RFC3339),
			Color:       0x0000ff,
			Description: "`!weather get [location]`",
		}

		// In case there is no location, provide message about it
		if msg, err := client.session.ChannelMessageSendEmbed(client.message.ChannelID, embeddedMsg); err != nil {
			log.Printf("error sending message | %v | %v", msg, err)
		}

		return
	}

	// Prepares an embedded message with weather information.
	embeddedMsg := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{},
		Title:       "Weather in " + res.Name + " | " + res.Weather[0].Main,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x0000ff,
		Description: "`!weather get [location]`",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Temperature",
				Value: fmt.Sprintf("%.1f °C", res.Main.Temp),
			},
			{
				Name:  "Full Weather Desc",
				Value: strings.Title(res.Weather[0].Description),
			},
			{
				Name:  "Humidity",
				Value: strconv.Itoa(res.Main.Humidity),
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http://openweathermap.org/img/w/" + res.Weather[0].Icon + ".png",
		},
	}

	// Sends the embedded message to the channel.
	if msg, err := client.session.ChannelMessageSendEmbed(client.message.ChannelID, embeddedMsg); err != nil {
		log.Printf("error sending message | %v | %v", msg, err)
	}
}

func (client *Client) Translate() {
	log.Printf("go.bot.handler.Translate request recieved")

	// Extracts the translation command from the message.
	cmd := client.ctx.Value("cmd").(string)
	cmd = strings.TrimPrefix(cmd, "translate ")
	client.ctx = context.WithValue(client.ctx, "cmd", cmd)

	// Defines translation-related routes and create a router instance.
	routes := map[*Route]interface{}{
		&Route{"to", "`!translate to [2-letters language code] [content]` - translates the text after to the give language"}: client.ToTranslate,
	}

	translateRouter := &Router{routes: routes}

	// Handles the translation command using the translation router.
	translateRouter.handleRoute(client)
}

func (client *Client) ToTranslate() {
	log.Printf("go.handler.Translate.To request recieved")

	// Extracts the translation target from the command.
	cmd := client.ctx.Value("cmd").(string)
	cmd = strings.TrimPrefix(cmd, "to ")
	client.ctx = context.WithValue(client.ctx, "cmd", cmd)

	log.Printf("%v", cmd)

	// Performs the translation asynchronously.
	go client.handleReaction(cmd)
}

func getLanguageCodeFromMessage(content string) (string, bool) {
	// This function parses the content to extract the language code.
	// It returns the language code and true if successful, or an empty string and false if not.

	// Splits the content into words
	words := strings.Fields(content)

	// Checks if there are at least two words (command and language code)
	if len(words) < 2 {
		return "", false // Not enough words in the message
	}

	// Extracts the language code from the second word
	languageCode := words[0]

	return languageCode, true
}

// handleReaction() handles reactions to translate messages.
func (client *Client) handleReaction(cmd string) {
	log.Printf("go.bot.Translate.To.handleReaction request recieved")
	// Extracts target language from message content
	targetLanguageCode, ok := getLanguageCodeFromMessage(cmd)
	if !ok {
		// Language code not provided or not recognized
		return
	}

	// Checks if the language code is valid
	target, ok := languages[targetLanguageCode]
	if !ok {
		// Not a valid language code
		embeddedMsg := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{},
			Title:       "Not a valid language code",
			Timestamp:   time.Now().Format(time.RFC3339),
			Color:       0x0000ff,
			Description: "`!translate to [2-letters language code] [content]`",
		}

		// In case there is no location, provide message about it
		if msg, err := client.session.ChannelMessageSendEmbed(client.message.ChannelID, embeddedMsg); err != nil {
			log.Printf("error sending message | %v | %v", msg, err)
		}
		return
	}
	target = targetLanguageCode

	// Trims the target language code from the command
	cmd = strings.TrimPrefix(cmd, target)
	client.ctx = context.WithValue(client.ctx, "cmd", cmd)
	content := cmd

	// Trims any leading whitespace from the content
	content = strings.TrimLeft(content, " ")

	// Translates the content to the target language
	source, text, err := translate(GoogleTranslateToken, target, content) //
	if err != nil {
		log.Printf("translate(%v, %q) error: %v", target, content, err)
		return
	}

	name := client.message.Author.Username

	// Creates an embed message with the translation
	embed := &discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: client.message.Author.AvatarURL("128"),
		},
		Description: text,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v (%v) → %v (%v)", languages[source], source, languages[target], target),
		},
	}

	// Sends the translation as an embed message to the channel
	if _, err := client.session.ChannelMessageSendEmbed(client.message.ChannelID, embed); err != nil {
		log.Printf("ChannelMessageSend(%q) error: %v", text, err)
		return
	}
}
