package internal

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"vrungel.maxvk.com/controller/config/discord"
)

// var (
// 	Session *discordgo.Session
// )

type DiscordBot struct {
	Session *discordgo.Session
}

func (bot *DiscordBot) Start(ctx context.Context) error {

	session, err := discordgo.New("Bot " + discord.Token)
	if err != nil {
		panic(err)
	}

	bot.Session = session

	bot.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	bot.Session.State.MaxMessageCount = 50

	err = bot.Session.Open()

	if err != nil {
		panic(err)
	}

	// Sends message to #general
	_, err = bot.Session.ChannelMessageSend("1392648744532443220", "Vrungel is starting...")
	if err != nil {
		fmt.Println("error sending message,", err)
	}

	<-ctx.Done()

	_, err = bot.Session.ChannelMessageSend("1392648744532443220", "Vrungel is stopping... Goodbye")
	if err != nil {
		fmt.Println("error sending message,", err)
	}

	return bot.Session.Close()
}
