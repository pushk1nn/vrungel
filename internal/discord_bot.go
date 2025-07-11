package internal

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"vrungel.maxvk.com/controller/config/discord"
)

// var (
// 	Session *discordgo.Session
// )

type DiscordBot struct {
	Session *discordgo.Session
}

func (bot *DiscordBot) Start(ctx context.Context) error {

	logger := log.FromContext(ctx)
	session, err := discordgo.New("Bot " + discord.Token)
	if err != nil {
		panic(err)
	}

	bot.Session = session
	bot.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	bot.Session.State.MaxMessageCount = 50

	err = bot.Session.Open()
	logger.Info("Started Bot")
	if err != nil {
		logger.Error(err, "Could not start Bot")
	}

	// Startup action
	_, err = bot.Session.ChannelMessageSend("1392648744532443220", "Vrungel is starting...")
	if err != nil {
		logger.Error(err, "Could not send startup message")
	}

	<-ctx.Done()
	logger.Info("Stopping Bot")

	// Shutdown action
	_, err = bot.Session.ChannelMessageSend("1392648744532443220", "Vrungel is stopping... Goodbye")
	if err != nil {
		logger.Error(err, "Could not send shutdown message")
	}

	return bot.Session.Close()
}
