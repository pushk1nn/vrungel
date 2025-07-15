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

	"github.com/bwmarrin/discordgo"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crdv1 "vrungel.maxvk.com/controller/api/crd/v1"
	"vrungel.maxvk.com/controller/internal/bot"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Setup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
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

	err = session.Open()
	if err != nil {
		logger.Error(err, "unable to open discord session")
	}

	r.BotManager.SetSession(session)

	logger.Info("discordgo session started successfully")
	// meta.SetStatusCondition(&setupList.Items[0].Status.Conditions, metav1.Condition{
	// 	Type:    "Ready",
	// 	Status:  metav1.ConditionTrue,
	// 	Reason:  "BotSessionEstablished",
	// 	Message: "Discord bot session established successfully",
	// })

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SetupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.Setup{}).
		Named("crd-setup").
		Complete(r)
}
