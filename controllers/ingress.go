package controllers

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "operators/WorkerBundle/api/v1"
)

func createIngressPaths(instance *apiv1.WorkerBundle) []networkingv1.HTTPIngressPath {
	paths := make([]networkingv1.HTTPIngressPath, len(instance.Spec.Workers))
	pathType := networkingv1.PathTypePrefix
	for i, worker := range instance.Spec.Workers {
		paths[i] = networkingv1.HTTPIngressPath{
			Path:     getIngressPathName(worker),
			PathType: &pathType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: getServiceName(worker.WorkerName),
					Port: networkingv1.ServiceBackendPort{
						Number: worker.WorkerNumber,
					},
				},
			},
		}
	}
	return paths
}

func createIngress(instance *apiv1.WorkerBundle) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getIngressName(instance.Spec.DeploymentName),
			Namespace: instance.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "worker.127.0.0.1.sslip.io",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: createIngressPaths(instance),
						},
					},
				},
			},
		},
	}
}
