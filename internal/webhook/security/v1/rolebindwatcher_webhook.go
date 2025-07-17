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

package v1

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	rbacv1 "k8s.io/api/rbac/v1"
)

// nolint:unused
// log is for logging in this package.
var rolebindwatcherlog = logf.Log.WithName("rolebindwatcher-resource")

// SetupRoleBindWatcherWebhookWithManager registers the webhook for RoleBindWatcher in the manager.
func SetupRoleBindWatcherWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&rbacv1.RoleBinding{}).
		WithValidator(&RoleBindWatcherCustomValidator{}).
		WithDefaulter(&RoleBindWatcherCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-security-vrungel-maxvk-com-v1-rolebindwatcher,mutating=true,failurePolicy=fail,sideEffects=None,groups=security.vrungel.maxvk.com,resources=rolebindwatchers,verbs=create;update,versions=v1,name=mrolebindwatcher-v1.kb.io,admissionReviewVersions=v1

// RoleBindWatcherCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind RoleBindWatcher when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type RoleBindWatcherCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &RoleBindWatcherCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind RoleBindWatcher.
func (d *RoleBindWatcherCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	rolebindwatcher, ok := obj.(*rbacv1.RoleBinding)

	if !ok {
		return fmt.Errorf("expected an RoleBindWatcher object but got %T", obj)
	}
	rolebindwatcherlog.Info("Defaulting for RoleBindWatcher", "name", rolebindwatcher.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-security-vrungel-maxvk-com-v1-rolebindwatcher,mutating=false,failurePolicy=fail,sideEffects=None,groups=security.vrungel.maxvk.com,resources=rolebindwatchers,verbs=create;update,versions=v1,name=vrolebindwatcher-v1.kb.io,admissionReviewVersions=v1

// RoleBindWatcherCustomValidator struct is responsible for validating the RoleBindWatcher resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type RoleBindWatcherCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &RoleBindWatcherCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type RoleBindWatcher.
func (v *RoleBindWatcherCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	rolebindwatcher, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("expected a RoleBindWatcher object but got %T", obj)
	}
	rolebindwatcherlog.Info("Validation for RoleBindWatcher upon creation", "name", rolebindwatcher.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, validateRoleBinding(rolebindwatcher)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type RoleBindWatcher.
func (v *RoleBindWatcherCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	rolebindwatcher, ok := newObj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("expected a RoleBindWatcher object for the newObj but got %T", newObj)
	}
	rolebindwatcherlog.Info("Validation for RoleBindWatcher upon update", "name", rolebindwatcher.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type RoleBindWatcher.
func (v *RoleBindWatcherCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	rolebindwatcher, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("expected a RoleBindWatcher object but got %T", obj)
	}
	rolebindwatcherlog.Info("Validation for RoleBindWatcher upon deletion", "name", rolebindwatcher.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

func validateRoleBinding(rb *rbacv1.RoleBinding) error {
	var allErrs field.ErrorList

	if err := validateRoleBindingRole(rb); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "security.vrungel.maxvk.com", Kind: "RoleBindWatcher"},
		rb.Name, allErrs)

}

func validateRoleBindingRole(rb *rbacv1.RoleBinding) *field.Error {

	if rb.RoleRef.Name == "dangerous-role" {
		return field.Invalid(field.NewPath("roleRef").Child("name"), rb.RoleRef.Name, "identity is being binded to a dangerous role")
	}
	return nil
}
