package handlers

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

type Constraint struct {
	Role string
}

func RoleConstraint(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Decode interaction data
	req := strings.SplitN(i.MessageComponentData().CustomID, ":", 2)[1]
	decoded, err := base64.StdEncoding.DecodeString(req)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return
	}

	// Call helper function to fill constraint struct with data
	constraint := populateConstraint(string(decoded))

	// Fill template
	var tmplFile = "roleconstraint.tmpl"
	var tmplPath = "templates/roleconstraint.tmpl"

	tmpl, err := template.New(tmplFile).ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, constraint)
	if err != nil {
		panic(err)
	}

	// Respond to the interaction
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Button clicked! base64 encoded request: %s", decoded),
		},
	})
	if err != nil {
		panic(err)
	}
}

func populateConstraint(req string) Constraint {
	body := strings.SplitN(req, "|", 3)

	return Constraint{
		Role: body[2],
	}
}
