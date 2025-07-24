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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	log "sigs.k8s.io/controller-runtime/pkg/log"

	securityv1 "vrungel.maxvk.com/controller/api/security/v1"
)

// RuleReconciler reconciles a Rule object
type RuleReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Initialized bool
	RiskyRoles  map[string]struct{}
}

// +kubebuilder:rbac:groups=security.vrungel.maxvk.com,resources=rules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.vrungel.maxvk.com,resources=rules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=security.vrungel.maxvk.com,resources=rules/finalizers,verbs=update

func (r *RuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var ruleList securityv1.RuleList

	if err := r.List(ctx, &ruleList); err != nil {
		panic(err)
	}

	roleSet := make(map[string]struct{})

	for _, instance := range ruleList.Items {
		for i := range instance.Spec.Risky {
			roleSet[instance.Spec.Risky[i]] = struct{}{}
		}
	}

	r.RiskyRoles = roleSet
	logger.Info("Ruleset initialized!")
	fmt.Println(r.RiskyRoles)

	r.Initialized = true

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.Rule{}).
		Named("security-rule").
		Complete(r)
}
