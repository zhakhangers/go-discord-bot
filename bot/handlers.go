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
	routes map[*Route]interface{}
}

type Route struct {
	name        string
	description string
}

type contextKey string

const (
	cmdKey    contextKey = "cmd"
	prefixKey contextKey = "!"
)

func (r *Router) handleRoute(c *Client) {
	cmd := c.ctx.Value("cmd").(string)

	if strings.HasPrefix(cmd, "help") {
		r.help(c)
	}

	for routeName, routeFunc := range r.routes {
		if strings.HasPrefix(cmd, routeName.name) {
			log.Printf("Route found based on command | %s -> %s", cmd, routeName)
			routeFunc.(func())()
		}
	}
}

func (r *Router) help(c *Client) {
	var res string

	for route := range r.routes {
		res = res + "\n" + route.description
	}

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

	if client.message.Author.ID == client.session.State.User.ID {
		return
	}

	if !strings.HasPrefix(client.message.Content, prefix) {
		return
	}

	log.Print("Prefix hit")
	cmd = strings.TrimPrefix(client.message.Content, prefix)
	client.ctx = context.WithValue(client.ctx, cmdKey, cmd)

	routes := map[*Route]interface{}{
		&Route{"weather", "`,weather help` - Weather info from Open Weather Maps "}: client.GetWeather,
	}

	router := &Router{routes: routes}

	router.handleRoute(&client)
}

func (client *Client) Weather() {
	log.Printf("go.bot.handler.Weather request recieved")

	cmd := client.ctx.Value("cmd").(string)
	cmd = strings.TrimPrefix(cmd, "weather ")
	client.ctx = context.WithValue(client.ctx, cmdKey, cmd)

	routes := map[*Route]interface{}{
		&Route{"get", "`,weather get [location]` - gets the weather for a location"}: client.GetWeather,
	}

	weatherRouter := &Router{routes: routes}

	weatherRouter.handleRoute(client)
}

func (client *Client) GetWeather() {
	log.Printf("go.bot.handler.Weather.Get request received")

	cmd := client.ctx.Value(cmdKey).(string)
	cmd = strings.TrimPrefix(cmd, "get ")
	client.ctx = context.WithValue(client.ctx, cmdKey, cmd)

	go client.processWeatherRequest(cmd)
}

func (client *Client) processWeatherRequest(cmd string) {
	weatherClient := NewWeatherClient(client.ctx)
	res, err := weatherClient.GetWeatherByLocation(cmd)
	if err != nil {
		log.Printf("error getting weather: %s | %v", cmd, err)
		return
	}

	embeddedMsg := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{},
		Title:       "Weather in " + res.Name + " | " + res.Weather[0].Main,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x0000ff,
		Description: "`,weather get [location]`",
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

	if msg, err := client.session.ChannelMessageSendEmbed(client.message.ChannelID, embeddedMsg); err != nil {
		log.Printf("error sending message | %v | %v", msg, err)
	}
}
