package bot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DiscordBotManager struct {
	mu      sync.RWMutex
	session *discordgo.Session
	Cache   *cache.Cache
}

func NewDiscordBotManager() *DiscordBotManager {
	return &DiscordBotManager{}
}

// Start is required for a "Runnable" to be registered with the controller-manager.
// Ironically, it is just setting up the bot to Stop with the rest of the program
// by waiting for the context to finish (program stopped). This way, it won't hang
// open after the program stops.
func (d *DiscordBotManager) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)

	<-ctx.Done()
	logger.Info("Stopping bot")

	if d.GetSession() != nil {
		return d.GetSession().Close()
	}
	return nil
}

func (d *DiscordBotManager) SetSession(s *discordgo.Session) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.session = s
}

func (d *DiscordBotManager) GetSession() *discordgo.Session {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.session
}

func (d *DiscordBotManager) DiscordLog(ctx context.Context, obj client.Object) *discordgo.Message {

	role, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		panic(ok)
	}

	message, err := d.GetSession().ChannelMessageSendComplex(
		"1393623353830412358",
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author:      &discordgo.MessageEmbedAuthor{},
					Color:       0xebad50, // Yellow
					Description: fmt.Sprintf("A %s has been detected in namespace %s", objType(obj), obj.GetNamespace()),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Role",
							Value:  role.RoleRef.Name,
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
					Title:     "Risky Role Detected",
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Block Role Binding",
							Style:    discordgo.PrimaryButton,
							CustomID: d.cacheEntry(obj),
						},
					},
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	return message
}

// Set k:v pair in cache and return the UUID to be used as the discord interaction CustomID
func (d *DiscordBotManager) cacheEntry(obj client.Object) string {
	id := fmt.Sprintf("role_constraint:%s", uuid.NewString()) // TODO: change this to be generalized

	// rb, ok := obj.(*rbacv1.RoleBinding) // TODO: Maybe this doesn't need to be asserted twice?
	// if !ok {
	// 	panic(ok)
	// }

	// // Data that will be logged to the cache
	// data := RoleBindingConstraint{
	// 	role: rb.RoleRef.Name,
	// }

	// Set cache entry
	d.Cache.Set(id, obj, cache.DefaultExpiration)

	return id
}

func objType(obj client.Object) string {

	switch obj.(type) {

	case *rbacv1.RoleBinding:
		return "RoleBinding"
	case *rbacv1.Role:
		return "Role"
	case *rbacv1.ClusterRole:
		return "ClusterRole"
	default:
		return "UnknownType"
	}
}
