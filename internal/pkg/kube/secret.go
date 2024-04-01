package kube

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-20 14:21:23
 * @file: secret.go
 * @description: secret operation
 */

// NewSecret returns a new secret
func NewSec(name, key string, req ctrl.Request) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-secret",
			Namespace: "default",
			Labels: map[string]string{
				"app": req.Name,
			},
		},
		Data: map[string][]byte{
			key: {},
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func NewSecret(name, key string, req ctrl.Request) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app": req.Name + "-secret",
			},
		},
		Data: map[string][]byte{
			key: {},
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func NewSecretEnv(envs []corev1.EnvVar, req ctrl.Request) *corev1.Secret {
	data := make(map[string][]byte)
	for _, env := range envs {
		data[env.Name] = []byte(env.ValueFrom.String())
	}

	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app": req.Name,
			},
		},
		Data: data,
		Type: corev1.SecretTypeOpaque,
	}
}

func NewSecretEnvFrom(envFromRefs []corev1.EnvFromSource, req ctrl.Request) *corev1.Secret {
	var secret *corev1.Secret // Declare the "secret" variable
	for _, envFrom := range envFromRefs {
		if envFrom.SecretRef != nil {
			secret = NewSecret(envFrom.SecretRef.Name, envFrom.SecretRef.Name, req)
			break
		}
	}
	return secret
}
