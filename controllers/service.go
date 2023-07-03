package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "operators/WorkerBundle/api/v1"
)

func createService(instance *apiv1.WorkerBundle) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getServiceName(instance.Spec.DeploymentName),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports:     []corev1.ServicePort{},
			Selector:  map[string]string{"app": getPodName(instance.Spec.DeploymentName)},
			ClusterIP: "None",
		},
	}
}
