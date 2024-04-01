package kube

import (
	singlev1 "github.com/keington/greatsql-operator/api/v1"
	"github.com/keington/greatsql-operator/internal/consts"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-18 17:06:11
 * @file: service.go
 * @description: kubenetes service operation
 */

func NewService(app *singlev1.Single) *corev1.Service {
	svcType := corev1.ServiceTypeClusterIP

	switch app.Spec.Type {
	case corev1.ServiceTypeClusterIP, corev1.ServiceTypeNodePort, corev1.ServiceTypeLoadBalancer:
		svcType = app.Spec.Type
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "Single",
				}),
			},
			Labels: map[string]string{
				consts.AppKubernetesComponent: "controller",
				consts.AppKubernetesName:      app.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  svcType,
			Ports: app.Spec.Ports,
			Selector: map[string]string{
				consts.AppKubernetesComponent: "controller",
				consts.AppKubernetesName:      app.Name,
			},
		},
	}
}
