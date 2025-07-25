package handlers

import (
	"fmt"
	"os"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	rbacv1 "k8s.io/api/rbac/v1"
	"vrungel.maxvk.com/controller/internal/bot/git"
	"vrungel.maxvk.com/controller/structs"
)

type HandlerManager struct {
	GitManager *git.GitManager
	Cache      *cache.Cache
}

func (h *HandlerManager) RoleConstraint(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Call helper function to fill constraint struct with data
	constraint := h.populateConstraint(i.MessageComponentData().CustomID)

	// Fill template
	var tmplFile = "roleconstraint.tmpl"
	var tmplPath = "templates/roleconstraint.tmpl"

	tmpl, err := template.New(tmplFile).ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}

	output := "/tmp/vrungel-automation/rolebinding-constraints/constraint.yaml"
	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	tmpl.Execute(f, constraint)
	f.Close()

	h.GitManager.Commit("constraint.yaml")
	h.GitManager.Push()

	// Respond to the interaction
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Processing constraint creation...",
		},
	})
	if err != nil {
		panic(err)
	}
}

func (h *HandlerManager) populateConstraint(id string) structs.Constraint {
	// Retrieve entry from request cache
	query, found := h.Cache.Get(id)
	if found {
		// Assert that type is role binding
		role := query.(*rbacv1.RoleBinding) // TODO: Generalize for different types

		return structs.Constraint{
			Target: role.RoleRef.Name,
		}
	}

	fmt.Println("Could not find CustomID in cache")
	return structs.Constraint{}
}
