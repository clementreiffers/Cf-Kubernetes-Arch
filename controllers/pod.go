package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "operators/WorkerBundle/api/v1"
)

func createPodPorts(ports []apiv1.Worker) []v1.ContainerPort {
	podPorts := make([]v1.ContainerPort, len(ports))
	for i, port := range ports {
		podPorts[i] = v1.ContainerPort{
			Name:          port.WorkerName,
			ContainerPort: port.WorkerNumber,
		}
	}
	return podPorts
}

func createPodSpec(instance *apiv1.WorkerBundle) v1.PodSpec {
	return v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  getPodName(instance.Spec.DeploymentName),
				Image: instance.Spec.PodTemplate.Image,
				Ports: createPodPorts(instance.Spec.Workers),
			},
		},
	}
}

func createDeployment(instance *apiv1.WorkerBundle) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Namespace: instance.GetNamespace(), Name: getDeploymentName(instance.Spec.DeploymentName)},
		Spec: appsv1.DeploymentSpec{
			Replicas: new(int32),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": getPodName(instance.Spec.DeploymentName)},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   getPodName(instance.Spec.DeploymentName),
					Labels: map[string]string{"app": getPodName(instance.Spec.DeploymentName)},
				},
				Spec: createPodSpec(instance),
			},
		},
	}
}
