package bot

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	ctx     context.Context
	session *discordgo.Session
	message *discordgo.MessageCreate
}
