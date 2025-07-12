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

package controller

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	bot "vrungel.maxvk.com/controller/internal/bot"
)

// PodTrackerReconciler reconciles a PodTracker object
type PodTrackerReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Discord bot.DiscordBot
}

// +kubebuilder:rbac:groups=crd.vrungel.maxvk.com,resources=podtrackers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crd.vrungel.maxvk.com,resources=podtrackers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=crd.vrungel.maxvk.com,resources=podtrackers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodTracker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *PodTrackerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Controller triggered")

	return ctrl.Result{}, nil
}

func (r *PodTrackerReconciler) HandleRBACEvents(ctx context.Context, rb client.Object) []reconcile.Request {

	logger := log.FromContext(ctx)

	if rb.GetNamespace() == "default" {
		logger.Info("Role binding found: " + rb.GetName())
		r.Discord.DiscordLog(rb)
	}
	return []reconcile.Request{}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodTrackerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("Vrungel").
		Watches(
			&rbacv1.RoleBinding{},
			handler.EnqueueRequestsFromMapFunc(r.HandleRBACEvents),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}
