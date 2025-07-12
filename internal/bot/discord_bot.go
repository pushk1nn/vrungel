package bot

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// var (
// 	Discord *discordgo.Discord
// )

type DiscordBot struct {
	Discord *discordgo.Session
}

func (bot *DiscordBot) Start(ctx context.Context) error {

	logger := log.FromContext(ctx)

	bot.Discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	bot.Discord.State.MaxMessageCount = 50

	err := bot.Discord.Open()
	logger.Info("Started Bot")
	if err != nil {
		logger.Error(err, "Could not start Bot")
	}

	// Startup action
	_, err = bot.Discord.ChannelMessageSend("1392648744532443220", "Vrungel is starting...")
	if err != nil {
		logger.Error(err, "Could not send startup message")
	}

	<-ctx.Done()
	logger.Info("Stopping Bot")

	// Shutdown action
	_, err = bot.Discord.ChannelMessageSend("1392648744532443220", "Vrungel is stopping... Goodbye")
	if err != nil {
		logger.Error(err, "Could not send shutdown message")
	}

	return bot.Discord.Close()
}

func (bot *DiscordBot) DiscordLog(objName string) {

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: "This is a discordgo embed",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Role Binding",
				Value:  objName,
				Inline: true,
			},
			{
				Name:   "I am a second field",
				Value:  "I am a value",
				Inline: true,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     "I am an Embed",
	}

	bot.Discord.ChannelMessageSendEmbed("1392648744532443220", embed)
}
