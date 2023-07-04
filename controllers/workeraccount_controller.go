/*
Copyright 2023 clementreiffers.

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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1 "operators/WorkerBundle/api/v1"
)

// WorkerAccountReconciler reconciles a WorkerAccount object
type WorkerAccountReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.cf-worker,resources=workeraccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.cf-worker,resources=workeraccounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.cf-worker,resources=workeraccounts/finalizers,verbs=update

func createWorkerBundle(instance *apiv1.WorkerAccount) apiv1.WorkerBundle {
	return apiv1.WorkerBundle{
		ObjectMeta: metav1.ObjectMeta{Name: instance.Spec.WorkerBundleName, Namespace: instance.GetNamespace()},
		Spec: apiv1.WorkerBundleSpec{
			DeploymentName: instance.Spec.WorkerBundleName,
			PodTemplate: apiv1.WorkerBundlePodTemplate{
				ImagePullSecret: instance.Spec.PodTemplate.ImagePullSecret,
				Image:           "nginx",
			},
		},
	}
}
func workerAccountApplyResource(r *WorkerAccountReconciler, ctx context.Context, resource client.Object, foundResource client.Object) error {
	err := r.Get(ctx, types.NamespacedName{Name: resource.GetName(), Namespace: resource.GetNamespace()}, foundResource)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(ctx, resource)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (r *WorkerAccountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.Log.WithValues("WorkerAccount", req.NamespacedName)

	instance := &apiv1.WorkerAccount{}
	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	workerBundle := createWorkerBundle(instance)
	err = workerAccountApplyResource(r, ctx, &workerBundle, &apiv1.WorkerBundle{})
	if err != nil {
		logger.Error(err, "unable to create WorkerBundle")
		return ctrl.Result{}, err
	}

	logger.Info("successfully created a worker bundle!")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkerAccountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.WorkerAccount{}).
		Complete(r)
}
