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

package security

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"vrungel.maxvk.com/controller/internal/bot"

	// securityv1 "vrungel.maxvk.com/controller/api/security/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	secv1 "vrungel.maxvk.com/controller/api/security/v1"
)

// RoleBindWatcherReconciler reconciles a RoleBindWatcher object
type RoleBindWatcherReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	BotManager  *bot.DiscordBotManager
	RuleManager *RuleReconciler
}

var (
	RoleBindWatcherList secv1.RoleBindWatcherList
)

// +kubebuilder:rbac:groups=security.vrungel.maxvk.com,resources=rolebindwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.vrungel.maxvk.com,resources=rolebindwatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=security.vrungel.maxvk.com,resources=rolebindwatchers/finalizers,verbs=update

func (r *RoleBindWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Requeue the Reconcile request if the ruleset from the Rule reconciler is not ready
	if !r.RuleManager.Initialized {
		logger.Info("RoleBinding ruleset not ready, requeuing reconcile request for: ")
		return ctrl.Result{RequeueAfter: time.Second * 3}, nil
	}

	if r.BotManager.GetSession() == nil {
		logger.Info("discordgo session not ready, requeuing reconcile request for: ")
		return ctrl.Result{RequeueAfter: time.Second * 3}, nil
	}

	// Client - get RoleBinding object
	var rb rbacv1.RoleBinding
	if err := r.Client.Get(ctx, req.NamespacedName, &rb); err != nil {
		logger.Info("could not find role binding object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if role is in risky role set
	if _, risky := r.RuleManager.RoleSet[rb.RoleRef.Name]; risky {
		r.BotManager.DiscordLog(ctx, &rb)
	}

	return ctrl.Result{}, nil
}

func (r *RoleBindWatcherReconciler) HandleRBACEvents(ctx context.Context, rb client.Object) []reconcile.Request {

	logger := log.FromContext(ctx)

	// r.generateLists(ctx)

	role, ok := rb.(*rbacv1.RoleBinding)
	if !ok {
		panic(ok)
	}

	logger.Info("RoleBinding found: " + role.GetName())
	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      rb.GetName(),
			Namespace: rb.GetNamespace(),
		},
	}}

}

// SetupWithManager sets up the controller with the Manager.
func (r *RoleBindWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(
			&rbacv1.RoleBinding{},
			handler.EnqueueRequestsFromMapFunc(r.HandleRBACEvents),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Named("security-rolebindwatcher").
		Complete(r)
}
