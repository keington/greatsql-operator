package kube

import (
	singlev1 "github.com/keington/greatsql-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-22 23:34:49
 * @file: pod.go
 * @description: kubernetes pod operation
 */

// NewContainers returns a new container
func NewContainers(app *singlev1.Single) []corev1.Container {
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range app.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	return []corev1.Container{
		{
			Name:            app.Name,
			Image:           app.Spec.PodSpec.Image,
			Resources:       app.Spec.PodSpec.Resources,
			StartupProbe:    &app.Spec.PodSpec.StartupProbe,
			ReadinessProbe:  &app.Spec.PodSpec.ReadinessProbe,
			LivenessProbe:   &app.Spec.PodSpec.LivenessProbe,
			SecurityContext: app.Spec.PodSpec.SecurityContext,
			Ports:           containerPorts,
			ImagePullPolicy: app.Spec.PodSpec.ImagePullPolicy,
			Env:             app.Spec.PodSpec.Envs,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      app.Name + "-config",
					MountPath: "/etc/my.cnf",
					SubPath:   "my.cnf",
				},
				{
					Name:      app.Name + "-db",
					MountPath: "/data",
				},
			},
		},
	}
}

// GetNodeName returns the node name of the pod
func GetNodeName(pod *corev1.Pod) string {
	if pod.Spec.NodeName != "" {
		return pod.Spec.NodeName
	}
	return ""
}
