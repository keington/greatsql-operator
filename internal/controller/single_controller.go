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

package controller

import (
	"context"
	"reflect"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/bytedance/sonic"
	singlev1 "github.com/keington/greatsql-operator/api/v1"
	"github.com/keington/greatsql-operator/internal/pkg/kube"
	"github.com/keington/greatsql-operator/internal/utils"
	"k8s.io/apimachinery/pkg/types"
)

// SingleReconciler reconciles a Single object
type SingleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	logger = ctrl.Log.WithName("greatsql-single-controller")
)

//+kubebuilder:rbac:groups=greatsql.greatsql.cn,resources=singles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=greatsql.greatsql.cn,resources=singles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=greatsql.greatsql.cn,resources=singles/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Modify the Reconcile function to compare the state specified by
// the Single object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *SingleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)

	log := logger.WithValues("Single GreatSql", req.NamespacedName)
	log.Info("Reconciling Single GreatSql...")

	singleGreatsql := &singlev1.Single{}
	// is the resource exists
	err := r.Client.Get(ctx, req.NamespacedName, singleGreatsql)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("SingleGreateSql resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch SingleGreateSql")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// handle finalizer
	finalizer := &utils.GreatSqlFinalizer{
		Cli:      r.Client,
		GreatSql: singleGreatsql,
	}
	if singleGreatsql.DeletionTimestamp != nil {
		if err := finalizer.HandleFinalizer(); err != nil {
			log.Error(err, "Could not handle finalizer")
			return ctrl.Result{}, err
		}

		finalizer.RemoveFinalizer()

		if err := r.Client.Update(ctx, singleGreatsql); err != nil {
			log.Error(err, "Could not update GreatSql")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// validate spec
	if err := r.validateSpec(singleGreatsql.Spec, req); err != nil {
		log.Error(err, "invalid spec, please check")
		return ctrl.Result{}, err
	}

	// if err := r.deleteAssociatedResources(ctx, req); err != nil {
	// 	log.Error(err, "Could not add finalizer")
	// 	return ctrl.Result{}, err

	// }

	if err := finalizer.AddFinalizer(); err != nil {
		log.Error(err, "Could not add finalizer")
		return ctrl.Result{}, err
	}

	// create deployment, persistentVolumeClaim and service
	deployGreatsql := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deployGreatsql); err != nil {

		// create configMap
		configMap := kube.NewConfigMap(req.Name+"-config", req.Namespace)
		//configMap.Annotations = map[string]string{consts.ConfigMapDataHash: kubernetes.GetConfigDataHash()}
		if err := r.Client.Create(ctx, configMap); err != nil {
			log.Error(err, "Could not create configMap")
			return ctrl.Result{}, err
		}
		log.Info("Create configMap is successful", "Name", configMap.Name, "Namespace", configMap.Namespace)

		pvc := kube.NewPersistentVolumeClaim(singleGreatsql)
		if err := r.Client.Create(ctx, pvc); err != nil {
			log.Error(err, "Could not create persistentVolumeClaim")
			return ctrl.Result{}, err
		}
		log.Info("Create persistentVolumeClaim is successful", "Name", pvc.Name, "Namespace", pvc.Namespace)

		//kubernetes.BindPersistentVolumeAndClaim(pv, pvc)

		// deployGreatsql not found, create it
		deploy := kube.NewDeployment(singleGreatsql, req.Name+"-config")
		// deploy := kubernetes.NewDeployment(singleGreatsql, configMap.Name)
		if err := r.Client.Create(ctx, deploy); err != nil {
			log.Error(err, "Could not create deployment")
			return ctrl.Result{}, err
		}
		log.Info("Create deployment is successful", "Name", deploy.Name, "Namespace", deploy.Namespace)

		// service not found, create it
		service := kube.NewService(singleGreatsql)
		if err := r.Client.Create(ctx, service); err != nil {
			log.Error(err, "Could not create service")
			return ctrl.Result{}, err
		}
		log.Info("Create service is successful", "Name", service.Name, "Namespace", service.Namespace)

		// Assign accessPoint value to singleGreatsql.Status.AccessPoint
		r.updateStatus(ctx, singleGreatsql, *service)
	}

	return r.watchResource(ctx, req, singleGreatsql)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SingleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&singlev1.Single{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

// setOwnerReference sets the owner reference of the persistentVolumeClaim and service
// TODO: SetControllerReference 在single中没有使用到，列为todo是为了在replica及mgr中使用，这两个涉及从属关系
// func (r *SingleReconciler) setOwnerReference(_ context.Context, singleGreatsql *singlev1.Single, persistentVolumeClaim *corev1.PersistentVolumeClaim, service *corev1.Service) error {
// 	if err := controllerutil.SetControllerReference(singleGreatsql, persistentVolumeClaim, r.Scheme); err != nil {
// 		return err
// 	}

// 	if err := controllerutil.SetControllerReference(singleGreatsql, service, r.Scheme); err != nil {
// 		return err
// 	}
// 	return nil
// }

// validateSpec validates the spec of the Single
func (r *SingleReconciler) validateSpec(spec singlev1.SingleSpec, req ctrl.Request) error {
	log := logger.WithValues("Request.Service.Namespace", req.Namespace, "Request.Service.Name", req.Name)

	// validate role and type

	// validate size
	if spec.Size == nil || *spec.Size == 0 {
		log.Error(nil, "size is required")
		return errors.NewBadRequest("size is required")
	}

	// validate podSpec
	if spec.PodSpec.Storage.PersistentVolumeClaimTemplate.StorageClassName == nil {
		log.Error(nil, "storageClassName is required")
		return errors.NewBadRequest("storageClassName is required")
	}

	return nil
}

// watchResource watches the resource
func (r *SingleReconciler) watchResource(ctx context.Context, req ctrl.Request, singleGreatsql *singlev1.Single) (ctrl.Result, error) {
	log := logger.WithValues("Request.Service.Namespace", req.Namespace, "Request.Service.Name", req.Name)

	// Association Annotations
	data, _ := sonic.Marshal(singleGreatsql.Spec)
	if singleGreatsql.Annotations != nil {
		singleGreatsql.Annotations["spec"] = string(data)
	} else {
		singleGreatsql.Annotations = map[string]string{"spec": string(data)}
	}

	if err := r.Client.Update(ctx, singleGreatsql); err != nil {
		log.Error(err, "Could not update GreatSql")
		return ctrl.Result{}, err
	}

	// Watch Resource
	oldSpec := &singlev1.SingleSpec{}
	if err := sonic.Unmarshal([]byte(singleGreatsql.Annotations["spec"]), oldSpec); err != nil {
		log.Error(err, "Could not unmarshal spec")
		return ctrl.Result{}, err
	}

	if reflect.DeepEqual(singleGreatsql.Spec, *oldSpec) {
		newDeployments := kube.NewDeployment(singleGreatsql, req.Name+"-config")
		oldDeployments := &appsv1.Deployment{}
		if err := r.Client.Get(ctx, req.NamespacedName, oldDeployments); err != nil {
			log.Error(err, "Could not get old deployment")
			return ctrl.Result{}, err
		}
		oldDeployments.Spec = newDeployments.Spec
		if err := r.Client.Update(ctx, oldDeployments); err != nil {
			log.Error(err, "Could not update deployment")
			return ctrl.Result{}, err
		}

		newResources := kube.NewService(singleGreatsql)
		oldService := &corev1.Service{}
		if err := r.Client.Get(ctx, req.NamespacedName, oldService); err != nil {
			log.Error(err, "Could not get old service")
			return ctrl.Result{}, err
		}
		oldService.Spec = newResources.Spec
		if err := r.Client.Update(ctx, oldService); err != nil {
			log.Error(err, "Could not update service")
			return ctrl.Result{}, nil
		}

		newConfigMap := kube.NewConfigMap(req.Name+"-config", req.Namespace)
		oldConfigMapName := req.Name + "-config"
		oldConfigMap := &corev1.ConfigMap{}
		if err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: oldConfigMapName}, oldConfigMap); err != nil {
			log.Error(err, "Could not get old configMap")
			return ctrl.Result{}, err
		}
		oldConfigMap.Data = newConfigMap.Data
		if err := r.Client.Update(ctx, oldConfigMap); err != nil {
			log.Error(err, "Could not update configMap")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// updateStatus updates the status of the Single
func (r *SingleReconciler) updateStatus(ctx context.Context, singleGreatsql *singlev1.Single, svc corev1.Service) error {
	log := logger.WithValues("Request.Service.Namespace", singleGreatsql.Namespace, "Request.Service.Name", singleGreatsql.Name)

	var accessPoint string
	switch svc.Spec.Type {
	case corev1.ServiceTypeClusterIP:
		accessPoint = svc.Spec.ClusterIP + ":" + strconv.Itoa(int(svc.Spec.Ports[0].Port))
	case corev1.ServiceTypeNodePort:
		accessPoint = svc.Spec.ClusterIP + ":" + strconv.Itoa(int(svc.Spec.Ports[0].NodePort))
	case corev1.ServiceTypeLoadBalancer:
		accessPoint = svc.Status.LoadBalancer.Ingress[0].IP + ":" + strconv.Itoa(int(svc.Spec.Ports[0].Port))
	}

	status := &singlev1.SingleStatus{
		AccessPoint: accessPoint,
		Size:        singleGreatsql.Spec.GetSize(),
		Ready:       0,
		Age:         svc.CreationTimestamp.String(),
	}

	if reflect.DeepEqual(singleGreatsql.Status, status) {
		return nil
	}

	singleGreatsql.Status = *status

	// update status
	if err := r.Client.Status().Update(ctx, singleGreatsql); err != nil {
		log.Error(err, "Could not update status")
		return err
	}

	return nil
}
