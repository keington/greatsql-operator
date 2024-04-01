package kube

import (
	singlev1 "github.com/keington/greatsql-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-18 22:30:23
 * @file: persistent_volume.go
 * @description: persistent volume claim
 */

var (
	volumeMode = corev1.PersistentVolumeFilesystem
)

// NewPersistentVolume returns a new persistent volume
func NewPersistentVolume(singleGreatsql *singlev1.Single) *corev1.PersistentVolume {

	persistentVolume := &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolume",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: singleGreatsql.Name + "-db",
		},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: *setDefaultStorage(singleGreatsql),
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRecycle,
			StorageClassName:              *singleGreatsql.Spec.PodSpec.Storage.PersistentVolumeClaimTemplate.StorageClassName,
			VolumeMode:                    &volumeMode,
			PersistentVolumeSource:        *singleGreatsql.Spec.PodSpec.Storage.PersistentVolumeSource,
		},
	}

	return persistentVolume
}

// NewPersistentVolumeClaim returns a new persistent volume claim
func NewPersistentVolumeClaim(singleGreatsql *singlev1.Single) *corev1.PersistentVolumeClaim {

	defaultStorage := setDefaultStorage(singleGreatsql)
	storageQuantity := defaultStorage.String()

	persistentVolumeClaim := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      singleGreatsql.Name + "-db",
			Namespace: singleGreatsql.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageQuantity),
				},
			},
			StorageClassName: singleGreatsql.Spec.PodSpec.Storage.PersistentVolumeClaimTemplate.StorageClassName,
			VolumeMode:       &volumeMode,
		},
	}

	return persistentVolumeClaim
}

// setDefaultStorage set default storage
func setDefaultStorage(singleGreatsql *singlev1.Single) *resource.Quantity {
	storageQuantity := resource.NewQuantity(singleGreatsql.Spec.PodSpec.Storage.PersistentVolumeClaimTemplate.Resources.Requests.Storage().Value(), resource.BinarySI)
	if storageQuantity == nil || storageQuantity.Value() <= 0 {
		storageQuantity = resource.NewQuantity(5, resource.BinarySI)
	}
	return storageQuantity
}
