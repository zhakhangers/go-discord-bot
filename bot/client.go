package bot

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Client represents a client instance that interacts with Discord.
type Client struct {
	ctx     context.Context          // Context for managing operations within the client.
	session *discordgo.Session       // Discord session associated with the client.
	message *discordgo.MessageCreate // The message being processed or responded to.
}
