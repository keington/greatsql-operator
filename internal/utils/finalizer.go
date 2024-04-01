package utils

import (
	"context"

	singlev1 "github.com/keington/greatsql-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

/**
 * @author: HuaiAn xu
 * @date: 2024-03-18 23:52:44
 * @file: finalizer.go
 * @description: resource finalizer
 */

type GreatSqlFinalizer struct {
	Cli      client.Client
	GreatSql *singlev1.Single
}

const (
	// greatSqlFinalizer is the finalizer name for the GreatSql
	greatSqlFinalizer = "finalizer.greatsql.cn"
)

var (
	logger = ctrl.Log.WithName("greatsql-finalizer")
)

func (g *GreatSqlFinalizer) HandleFinalizer() error {
	logger.WithValues("Request.Finalizer.Namespace", g.GreatSql.Namespace, "Request.Finalizer.Name", g.GreatSql.Name)

	if g.GreatSql.ObjectMeta.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(g.GreatSql, greatSqlFinalizer) {
			if err := g.finalizelPersistentVolumeClaim(); err != nil {
				return err
			}

			if err := g.finalizeDeployment(); err != nil {
				return err
			}

			if err := g.finalizerConfigMap(); err != nil {
				return err
			}
			controllerutil.RemoveFinalizer(g.GreatSql, greatSqlFinalizer)
			if err := g.Cli.Update(context.Background(), g.GreatSql); err != nil {
				logger.Error(err, "Could not remove finalizer from GreatSql")
				return err
			}
		}
	}

	return nil
}

// AddFinalizer adds the finalizer to the GreatSql
func (g *GreatSqlFinalizer) AddFinalizer() error {
	logger.WithValues("Request.Finalizer.Namespace", g.GreatSql.Namespace, "Request.Finalizer.Name", g.GreatSql.Name)

	if !controllerutil.ContainsFinalizer(g.GreatSql, greatSqlFinalizer) {
		controllerutil.AddFinalizer(g.GreatSql, greatSqlFinalizer)
		if err := g.Cli.Update(context.Background(), g.GreatSql); err != nil {
			logger.Error(err, "Could not add finalizer to GreatSql")
			return err
		}
	}

	return nil
}

// RemoveFinalizer removes the finalizer from the GreatSql
func (g *GreatSqlFinalizer) RemoveFinalizer() error {
	logger.WithValues("Request.Finalizer.Namespace", g.GreatSql.Namespace, "Request.Finalizer.Name", g.GreatSql.Name)

	if controllerutil.ContainsFinalizer(g.GreatSql, greatSqlFinalizer) {
		controllerutil.RemoveFinalizer(g.GreatSql, greatSqlFinalizer)
		if err := g.Cli.Update(context.Background(), g.GreatSql); err != nil {
			logger.Error(err, "Could not remove finalizer from GreatSql")
			return err
		}
	}

	return nil
}

// finalizelPersistentVolumeClaim removes the PVCs
func (g *GreatSqlFinalizer) finalizelPersistentVolumeClaim() error {
	logger.WithValues("Request.Finalizer.Namespace", g.GreatSql.Namespace, "Request.Finalizer.Name", g.GreatSql.Name)

	for i := 0; i < int(g.GreatSql.Spec.GetSize()); i++ {
		// pvcName := g.GreatSql.Name + "-" + g.GreatSql.Name + "-" + strconv.Itoa(i)
		pvcName := g.GreatSql.Name + "-db"
		// delete pvc
		err := g.Cli.Delete(context.TODO(), &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName,
				Namespace: g.GreatSql.Namespace,
			},
		})
		if err != nil && !errors.IsNotFound(err) {
			logger.Error(err, "Could not delete PersistentVolumeClaim "+pvcName)
			return err
		}
	}

	return nil
}

func (g *GreatSqlFinalizer) finalizeDeployment() error {
	logger.WithValues("Request.Finalizer.Namespace", g.GreatSql.Namespace, "Request.Finalizer.Name", g.GreatSql.Name)

	// delete deployment
	err := g.Cli.Delete(context.TODO(), &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.GreatSql.Name,
			Namespace: g.GreatSql.Namespace,
		},
	})
	if err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "Could not delete Deployment "+g.GreatSql.Name)
		return err
	}

	return nil
}

func (g *GreatSqlFinalizer) finalizerConfigMap() error {
	logger.WithValues("Request.Finalizer.Namespace", g.GreatSql.Namespace, "Request.Finalizer.Name", g.GreatSql.Name)

	// delete configmap
	err := g.Cli.Delete(context.TODO(), &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.GreatSql.Name + "config",
			Namespace: g.GreatSql.Namespace,
		},
	})
	if err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "Could not delete ConfigMap "+g.GreatSql.Name)
		return err
	}

	return nil
}
