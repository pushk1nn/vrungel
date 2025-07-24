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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleBindWatcherSpec defines the desired state of RoleBindWatcher.
type RoleBindWatcherSpec struct {
	Risky []string `json:"risky,omitempty"`
}

// RoleBindWatcherStatus defines the observed state of RoleBindWatcher.
type RoleBindWatcherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RoleBindWatcher is the Schema for the rolebindwatchers API.
type RoleBindWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleBindWatcherSpec   `json:"spec,omitempty"`
	Status RoleBindWatcherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoleBindWatcherList contains a list of RoleBindWatcher.
type RoleBindWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleBindWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RoleBindWatcher{}, &RoleBindWatcherList{})
}
