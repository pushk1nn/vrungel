package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RoleConstraint(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Respond to the interaction
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Button clicked! base64 encoded request: %s", i.MessageComponentData().CustomID),
		},
	})
	if err != nil {
		panic(err)
	}
}
