package kube

import (
	singlev1 "github.com/keington/greatsql-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-18 22:43:46
 * @file: statefulset.go
 * @description: statefulset operation
 */

func NewStatefulSet(singleGreatsql *singlev1.Single) *appsv1.StatefulSet {

	statefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      singleGreatsql.Name + "-statefulset",
			Namespace: singleGreatsql.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: singleGreatsql.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": singleGreatsql.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": singleGreatsql.Name,
					},
				},
			},
			// TODO: add volumeClaimTemplates
		},
	}

	return statefulSet
}
