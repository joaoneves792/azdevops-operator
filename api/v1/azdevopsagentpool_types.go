/*
Copyright 2023.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AzDevopsAgentPoolAutoscaling defines the limits and thresholds for autoscaling the pool
type AzDevopsAgentPoolAutoscaling struct {
	Max int32 `json:"max,omitempty"`
	Min int32 `json:"min,omitempty"`
}

type AzDevopsProject struct {
	Url          string `json:"url,omitempty"`
	PoolName     string `json:"poolName,omitempty"`
	ProjectName  string `json:"projectName,omitempty"`
	PatSecretRef string `json:"PATSecretRef"`
}

// AzDevopsAgentPoolSpec defines the desired state of AzDevopsAgentPool
type AzDevopsAgentPoolSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AzDevopsAgentPool. Edit azdevopsagentpool_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	Project     AzDevopsProject              `json:"project,omitempty"`
	Autoscaling AzDevopsAgentPoolAutoscaling `json:"autoscaling,omitempty"`
	Template    corev1.PodTemplateSpec       `json:"template,omitempty"`
}

// AzDevopsAgentPoolStatus defines the observed state of AzDevopsAgentPool
type AzDevopsAgentPoolStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AzDevopsAgentPool is the Schema for the azdevopsagentpools API
type AzDevopsAgentPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzDevopsAgentPoolSpec   `json:"spec,omitempty"`
	Status AzDevopsAgentPoolStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AzDevopsAgentPoolList contains a list of AzDevopsAgentPool
type AzDevopsAgentPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzDevopsAgentPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzDevopsAgentPool{}, &AzDevopsAgentPoolList{})
}
