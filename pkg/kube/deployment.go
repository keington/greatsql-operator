package kube

import (
	singlev1 "github.com/keington/greatsql-operator/api/v1"
	"github.com/keington/greatsql-operator/internal/consts"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-18 18:02:46
 * @file: deployment.go
 * @description: kubernetes deployment operation
 */

// NewDeployment returns a new deployment
func NewDeployment(single *singlev1.Single, configMapName string) *appsv1.Deployment {
	labels := map[string]string{
		consts.AppKubernetesComponent: "controller",
		consts.AppKubernetesName:      single.Name,
	}
	selector := &metav1.LabelSelector{MatchLabels: labels}

	affinity := setAffinity(single, labels)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      single.Name,
			Namespace: single.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(single, schema.GroupVersionKind{
					Group:   appsv1.SchemeGroupVersion.Group,
					Version: appsv1.SchemeGroupVersion.Version,
					Kind:    "Single",
				}),
			},
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: single.Spec.Size,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers:                    NewContainers(single),
					TerminationGracePeriodSeconds: single.Spec.PodSpec.TerminationGracePeriodSeconds,
					SchedulerName:                 single.Spec.PodSpec.SchedulerName,
					Affinity:                      affinity,
					ServiceAccountName:            single.Spec.PodSpec.ServiceAccountName,
					SecurityContext:               single.Spec.PodSpec.PodSecurityContext,
					NodeSelector:                  single.Spec.PodSpec.NodeSelector,
					Tolerations:                   single.Spec.PodSpec.Tolerations,
					Volumes: []corev1.Volume{
						{
							Name: single.Name + "-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMapName,
									},
									DefaultMode: &[]int32{0664}[0],
								},
							},
						},
						{
							Name: single.Name + "-db",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: single.Name + "-db",
								},
							},
						},
					},
					DNSPolicy: single.Spec.DnsPolicy,
				},
			},
			Selector: selector,
			Strategy: appsv1.DeploymentStrategy{
				Type: single.Spec.UpdateStrategy,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{IntVal: 1},
					MaxSurge:       &intstr.IntOrString{IntVal: 1},
				},
			},
		},
	}
}

// setAffinity set affinity and anti-affinity
func setAffinity(single *singlev1.Single, labels map[string]string) *corev1.Affinity {
	if single.Spec.PodSpec.Affinity.Advanced != nil {
		return nil
	}

	return &corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					TopologyKey: *single.Spec.PodSpec.Affinity.TopologyKey,
				},
			},
		},
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					TopologyKey: *single.Spec.PodSpec.Affinity.TopologyKey,
				},
			},
		},
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      *single.Spec.PodSpec.Affinity.TopologyKey,
								Operator: corev1.NodeSelectorOpNotIn,
								Values:   []string{""},
							},
						},
					},
				},
			},
		},
	}
}
