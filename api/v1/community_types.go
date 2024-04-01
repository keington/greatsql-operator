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
	corev1 "k8s.io/api/core/v1"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-17 18:34:07
 * @file: community_types.go
 * @description: common types
 */

type GretaSql struct {
	SingleSpec
}

// GreatSqlType defines the type of the GreatSql
// Supported values are "Single" "ReplicaofCluster" "SinglePrimaryGroupCluster" "MultiPrimaryGroupCluster"
// Single: Single instance of GreatSql
// TODO: ReplicaofCluster: replica of a GreatSql cluster(Replicaof)
// TODO: SinglePrimaryGroupCluster: group of GreatSql clusters(SinglePrimaryMGR)
// TODO: MultiPrimaryGroupCluster: group of GreatSql clusters(MultiPrimaryMGR)
type GreatSqlType string

const (
	GreatSqlTypeSingle                    GreatSqlType = "single"
	GreatSqlTypeReplicaofCluster          GreatSqlType = "replicaofCluster"
	GreatSqlTypeSinglePrimaryGroupCluster GreatSqlType = "singlePrimaryGroupCluster"
	GreatSqlTypeMultiPrimaryGroupCluster  GreatSqlType = "multiPrimaryGroupCluster"
)

type MemberRole string

const (
	SingleRole     MemberRole = "single"
	PrimaryRole    MemberRole = "primary"
	SencondaryRole MemberRole = "sencondary"
	ReplicaofRole  MemberRole = "replicaof"
)

// PodSpec defines the desired state of Pod
type PodSpec struct {
	Affinity                      *PodAffinity               `json:"affinity,omitempty"` // pod affinity(pod亲和性)
	Annotation                    map[string]string          `json:"annotation,omitempty"`
	Labels                        map[string]string          `json:"labels,omitempty"`
	NodeSelector                  map[string]string          `json:"nodeSelector,omitempty"`
	VolumeSpec                    *VolumeSpec                `json:"volumeSpec,omitempty"`
	Tolerations                   []corev1.Toleration        `json:"tolerations,omitempty"`                   //schedule tolerations
	TerminationGracePeriodSeconds *int64                     `json:"terminationGracePeriodSeconds,omitempty"` // 在规定时间内停止pod，俗称 优雅停机
	SchedulerName                 string                     `json:"schedulerName,omitempty"`
	PodSecurityContext            *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	ServiceAccountName            string                     `json:"serviceAccountName,omitempty"`
	Version                       string                     `json:"version,omitempty"`
	//Configurations                map[string]string          `json:"configurations,omitempty"`
	ContainerSpec `json:",inline"` // container spec
	Storage       *Storage         `json:"storage,omitempty"`
}

type Storage struct {
	Type                          string                            `json:"type,omitempty"`
	PersistentVolumeSource        *corev1.PersistentVolumeSource    `json:"persistentVolumeSource,omitempty"`
	PersistentVolumeClaimTemplate *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimTemplate,omitempty"`
}

// PodAffinity defines the affinity/anti-affinity rules for the pod.
type PodAffinity struct {
	// TODO:antiAffinityTopologyKey pod anti-affinity parameters
	//+builder:default="kubernetes.io/hostname"
	//+Optional
	TopologyKey *string          `json:"antiAffinityTopologyKey,omitempty"`
	Advanced    *corev1.Affinity `json:"advanced,omitempty"`
}

// VolumeSpec defines the volume spec for mysql.
type VolumeSpec struct {
	// EmptyDir to use as data volume for mysql. EmptyDir represents a temporary
	// directory that shares a pod's lifetime.
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`

	//  HostPath to use as data volume for mysql. HostPath represents a
	// pre-existing file or directory on the host machine that is directly
	// exposed to the container.
	HostPath *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`

	// PersistentVolumeClaim to specify PVC spec for the volume for mysql data.
	// It has the highest level of precedence, followed by HostPath and EmptyDir.
	// And represents the PVC specification.
	PersistentVolumeClaim *corev1.PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
}

// ContainerSpec defines the desired state of the container
type ContainerSpec struct {
	Image            string                        `json:"image"`                      // Image of the container
	ImagePullPolicy  corev1.PullPolicy             `json:"imagePullPolicy,omitempty"`  // Image pull policy
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"` // Image pull secrets
	Resources        corev1.ResourceRequirements   `json:"resources,omitempty"`        // Resource requirements
	StartupProbe     corev1.Probe                  `json:"startupProbe,omitempty"`     // Startup probe
	ReadinessProbe   corev1.Probe                  `json:"readinessProbe,omitempty"`   // Readiness probe
	LivenessProbe    corev1.Probe                  `json:"livenessProbe,omitempty"`    // Liveness probe
	SecurityContext  *corev1.SecurityContext       `json:"securityContext,omitempty"`  // Security context for the container
	Envs             []corev1.EnvVar               `json:"envs,omitempty"`             // Environment variables
}

// UpgradeOptions defines the desired state of UpgradeOptions
type UpgradeOptions struct {
	VersionServiceEndpoint string `json:"versionServiceEndpoint,omitempty"`
	Apply                  string `json:"apply,omitempty"`
}
