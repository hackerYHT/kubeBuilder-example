/*
Copyright 2024 yht.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodsbookSpec defines the desired state of Podsbook
type PodsbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image,Replica is an example field of Podsbook. Edit podsbook_types.go to remove/update
	Image   *string `json:"image,omitempty"`
	Replica *int32  `json:"replica,omitempty"`
}

// PodsbookStatus defines the observed state of Podsbook
type PodsbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	RealReplica int32 `json:"realReplica,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.realReplica",name=RealReplica,type=integer

// Podsbook is the Schema for the podsbooks API
type Podsbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodsbookSpec   `json:"spec,omitempty"`
	Status PodsbookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodsbookList contains a list of Podsbook
type PodsbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Podsbook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Podsbook{}, &PodsbookList{})
}
