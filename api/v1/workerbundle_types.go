/*
Copyright 2023 clementreiffers.

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

type Worker struct {
	WorkerName   string `json:"workerName"`
	WorkerNumber int32  `json:"workerNumber"`
	EnvPrefix    string `json:"envPrefix"`
	SecretRef    string `json:"secretRef"`
}

type WorkerBundlePodTemplate struct {
	Image           string `json:"image,omitempty"`
	ImagePullSecret string `json:"imagePullSecret"`
}

// WorkerBundleSpec defines the desired state of WorkerBundle
type WorkerBundleSpec struct {
	DeploymentName string                  `json:"deploymentName"`
	Workers        []Worker                `json:"workers,omitempty"`
	PodTemplate    WorkerBundlePodTemplate `json:"podTemplate"`
}

// WorkerBundleStatus defines the observed state of WorkerBundle
type WorkerBundleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// WorkerBundle is the Schema for the workerbundles API
type WorkerBundle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkerBundleSpec   `json:"spec,omitempty"`
	Status WorkerBundleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WorkerBundleList contains a list of WorkerBundle
type WorkerBundleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkerBundle `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkerBundle{}, &WorkerBundleList{})
}
