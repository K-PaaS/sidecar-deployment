/*
Copyright 2022.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TaskWorkloadSpec defines the desired state of TaskWorkload
type TaskWorkloadSpec struct {
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Required
	Command []string `json:"command"`

	// +kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets"`

	// +kubebuilder:validation:Optional
	Env []corev1.EnvVar `json:"env"`
}

// TaskWorkloadStatus defines the observed state of TaskWorkload
type TaskWorkloadStatus struct {
	//+kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration captures the latest generation of the TaskWorkload that has been reconciled
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TaskWorkload is the Schema for the taskworkloads API
type TaskWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskWorkloadSpec   `json:"spec,omitempty"`
	Status TaskWorkloadStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TaskWorkloadList contains a list of TaskWorkload
type TaskWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskWorkload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TaskWorkload{}, &TaskWorkloadList{})
}

func (t TaskWorkload) StatusConditions() []metav1.Condition {
	return t.Status.Conditions
}
