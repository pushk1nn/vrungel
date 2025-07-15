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

type Reporter struct {
	Kind    string `json:"kind,omitempty"`
	Key     string `json:"key,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// SetupSpec defines the desired state of Setup.
type SetupSpec struct {
	Name   string   `json:"name,omitempty"`
	Report Reporter `json:"report,omitempty"`
}

// SetupStatus defines the observed state of Setup.
type SetupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Setup is the Schema for the setups API.
type Setup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SetupSpec   `json:"spec,omitempty"`
	Status SetupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SetupList contains a list of Setup.
type SetupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Setup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Setup{}, &SetupList{})
}
