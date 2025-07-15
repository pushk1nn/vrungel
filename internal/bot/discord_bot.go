package bot

import (
	// "context"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/log"
)

type DiscordBotManager struct {
	mu      sync.RWMutex
	session *discordgo.Session
}

func NewDiscordBotManager() *DiscordBotManager {
	return &DiscordBotManager{}
}

// func (d *DiscordBotManager) Start(ctx context.Context) error {
// 	logger := log.FromContext(ctx)

// 	<-ctx.Done()
// 	logger.Info("Stopping bot")

// 	return d.session.Close()
// }

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

func (d *DiscordBotManager) DiscordLog(obj client.Object) {

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: fmt.Sprintf("A %s has been created in namespace %s", objType(obj), obj.GetNamespace()),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Kind",
				Value:  obj.GetName(),
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
		Title:     "Resource Creation",
	}

	d.GetSession().ChannelMessageSendEmbed("1393623353830412358", embed)
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
