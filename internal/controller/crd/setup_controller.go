/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package crd

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	gogit "github.com/go-git/go-git/v6"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crdv1 "vrungel.maxvk.com/controller/api/crd/v1"
	"vrungel.maxvk.com/controller/internal/bot"
	"vrungel.maxvk.com/controller/internal/bot/git"
	"vrungel.maxvk.com/controller/internal/bot/handlers"
)

// SetupReconciler reconciles a Setup object
type SetupReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	BotManager *bot.DiscordBotManager
}

// +kubebuilder:rbac:groups=crd.vrungel.maxvk.com,resources=setups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crd.vrungel.maxvk.com,resources=setups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=crd.vrungel.maxvk.com,resources=setups/finalizers,verbs=update

func (r *SetupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting discordgo session...")

	var setupList crdv1.SetupList

	if err := r.List(ctx, &setupList); err != nil {
		logger.Error(err, "unable to fetch setup-controller list")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	token := setupList.Items[0].Spec.Report.Key

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error(err, "unable to start discord session")
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	session.State.MaxMessageCount = 50

	err = session.Open()
	if err != nil {
		logger.Error(err, "unable to start discordgo session")
	}
	logger.Info("started discordgo session")

	g := InitGitManager()

	session.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Type == discordgo.InteractionMessageComponent {

				// prefix will be the name of the handler without the base64 encoded info used for handling the request
				prefix := strings.SplitN(i.MessageComponentData().CustomID, ":", 2)[0]

				h := &handlers.HandlerManager{
					GitManager: g,
				}

				switch prefix {
				case "role_constraint":
					h.RoleConstraint(s, i)
				}
			}
		},
	)

	r.BotManager.SetSession(session)

	logger.Info("discordgo session created successfully")

	return ctrl.Result{}, nil
}

func InitGitManager() *git.GitManager {
	path := "/tmp/vrungel-automation"

	r, err := gogit.PlainClone(path, &gogit.CloneOptions{
		URL:               "https://github.com/pushk1nn/argocd-test.git",
		RecurseSubmodules: gogit.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		panic(err)
	}

	return &git.GitManager{
		Path: path,
		Repo: r,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SetupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.Setup{}).
		Named("crd-setup").
		Complete(r)
}
