package internal

import (
	"github.com/bwmarrin/discordgo"
	"vrungel.maxvk.com/controller/config/discord"
)

var (
	Session *discordgo.Session
)

func init() {

	session, err := discordgo.New("Bot " + discord.Token)
	if err != nil {
		panic(err)
	}

	Session = session

	Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	Session.State.MaxMessageCount = 50
}
