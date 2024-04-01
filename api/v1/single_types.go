/*
Copyright 2024 greatsql.

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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-17 18:32:59
 * @file: single_types.go
 * @description: single types
 */

// SingleSpec defines the desired state of Single
type SingleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+kubebuilder:validation:Enum=single;replicaofGroupCluster;singlePrimaryGroupCluster;multiPrimaryGroupCluster
	GreatSqlType   GreatSqlType                  `json:"greatSqlType,omitempty"`
	Role           MemberRole                    `json:"role,omitempty"`
	Size           *int32                        `json:"size,omitempty"`
	PodSpec        PodSpec                       `json:"podSpec,omitempty"`
	Ports          []corev1.ServicePort          `json:"ports,omitempty"`
	Type           corev1.ServiceType            `json:"type,omitempty"`
	DnsPolicy      corev1.DNSPolicy              `json:"dnsPolicy,omitempty"`
	UpgradeOptions UpgradeOptions                `json:"upgradeOptions,omitempty"`
	UpdateStrategy appsv1.DeploymentStrategyType `json:"updateStrategy,omitempty"`
}

// GetSize returns the size of the single
func (s *SingleSpec) GetSize() int32 {
	if s.Size != nil {
		return *s.Size
	}
	return 1
}

// SingleStatus defines the observed state of Single
type SingleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AccessPoint             string `json:"accessPoint,omitempty"`
	Size                    int32  `json:"size,omitempty"`
	Ready                   int32  `json:"ready,omitempty"`
	Age                     string `json:"age,omitempty"`
	appsv1.DeploymentStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="AccessPoint",type="string",JSONPath=".status.accessPoint",description="The access point of the single"
//+kubebuilder:printcolumn:name="Size",type="integer",JSONPath=".spec.size",description="The size of the single"
//+kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.ready",description="The ready of the single"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".status.age",description="The age of the single"

// Single is the Schema for the singles API
type Single struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SingleSpec   `json:"spec,omitempty"`
	Status SingleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SingleList contains a list of Single
type SingleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Single `json:"items"`
}

func (s *Single) Finalizer() []string {
	return []string{"finalizer.single.greatsql.cn"}
}

func init() {
	SchemeBuilder.Register(&Single{}, &SingleList{})
}
