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
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1 "operators/WorkerBundle/api/v1"
)

// WorkerReleaseReconciler reconciles a WorkerRelease object
type WorkerReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.cf-worker,resources=workerreleases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.cf-worker,resources=workerreleases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.cf-worker,resources=workerreleases/finalizers,verbs=update

func getAllScriptsUrls(instance *apiv1.WorkerRelease) []string {
	values := make([]string, 0, len(instance.Spec.WorkerVersions))
	for _, value := range instance.Spec.WorkerVersions {
		values = append(values, value)
	}
	return values
}

func getAllScriptNames(instance *apiv1.WorkerRelease) []string {
	keys := make([]string, 0, len(instance.Spec.WorkerVersions))
	for key := range instance.Spec.WorkerVersions {
		keys = append(keys, key)
	}
	return keys
}

func createJobBuilder(instance *apiv1.WorkerRelease, bundleName string) apiv1.JobBuilder {
	return apiv1.JobBuilder{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.Accounts,
			Namespace: instance.GetNamespace()},
		Spec: apiv1.JobBuilderSpec{
			ScriptUrls:       getAllScriptsUrls(instance),
			TargetImage:      fmt.Sprintf("clementreiffers/build-%s", instance.Spec.Accounts),
			WorkerBundleName: bundleName,
			ScriptNames:      getAllScriptNames(instance),
		},
	}
}

func (r *WorkerReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.Log.WithValues("JobBuilder", req.NamespacedName)

	instance := &apiv1.WorkerRelease{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	workerAccount := apiv1.WorkerAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Spec.Accounts, Namespace: instance.GetNamespace()}, &workerAccount)
	if err != nil {
		return ctrl.Result{}, err
	}

	bundleName := workerAccount.Spec.WorkerBundleName

	jobBuilder := apiv1.JobBuilder{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Spec.Accounts, Namespace: instance.GetNamespace()}, &jobBuilder)
	if err != nil {
		jobBuilder := createJobBuilder(instance, bundleName)
		err = r.Create(ctx, &jobBuilder)
		if err != nil {
			logger.Error(err, "unable to create a JobBuilder")
			return ctrl.Result{}, err
		}
		logger.Info("JobBuilder created!")
		return ctrl.Result{}, nil
	} else {
		err = r.Delete(ctx, &jobBuilder)
		if err != nil {
			logger.Error(err, "unable to destroy the job builder")
			return ctrl.Result{}, err
		}
		jobBuilder := createJobBuilder(instance, bundleName)
		err = r.Create(ctx, &jobBuilder)
		if err != nil {
			logger.Error(err, "unable to create a JobBuilder")
			return ctrl.Result{}, err
		}
		logger.Info("JobBuilder created!")
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkerReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.WorkerRelease{}).
		Complete(r)
}
