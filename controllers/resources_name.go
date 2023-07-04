package controllers

import apiv1 "operators/WorkerBundle/api/v1"

func getPodName(instance string) string {
	return instance + "-pod"
}

func getServiceName(instance string) string {
	return instance + "-svc"
}

func getIngressName(instance string) string {
	return instance + "-ingress"
}

func getIngressPathName(port apiv1.Worker) string {
	return "/" + port.WorkerName
}

func getDeploymentName(instance string) string {
	return instance + "-depl"
}

func getJobName(instance string) string {
	return instance + "-job"
}
