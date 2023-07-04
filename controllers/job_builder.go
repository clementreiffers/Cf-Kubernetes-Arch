package controllers

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "operators/WorkerBundle/api/v1"
	"strings"
)

func generateAwsConfig() []v1.EnvVar {
	return []v1.EnvVar{
		{Name: "AWS_PROFILE", Value: "default"},
		{Name: "AWS_ENDPOINT", Value: "https://s3.fr-par.scw.cloud"},
		{Name: "AWS_BUCKET", Value: "stage-cf-worker"},
	}
}

func generateAwsCommandSync(scriptUrls []string, finalPath string) string {
	personalizedAws := "aws --endpoint-url=$(AWS_ENDPOINT) s3 sync"
	var commands []string
	for _, url := range scriptUrls {
		commands = append(commands, fmt.Sprintf(" %s %s %s", personalizedAws, url, finalPath))
	}
	return strings.Join(commands, " && ")
}

func generateDownloadFilesContainer(instance *apiv1.JobBuilder) v1.Container {
	return v1.Container{
		Name:            "download-files",
		Image:           "public.ecr.aws/aws-cli/aws-cli:latest",
		ImagePullPolicy: "IfNotPresent",
		Env:             generateAwsConfig(),
		VolumeMounts: []v1.VolumeMount{
			{Name: "s3-config", MountPath: "/root/.aws", ReadOnly: true},
			{Name: "context", MountPath: "/context"},
		},
		Command: []string{"/bin/sh"},
		Args:    []string{"-c", generateAwsCommandSync(instance.Spec.ScriptUrls, "/context")},
	}
}

func generateGettingDockerfile() v1.Container {
	return v1.Container{
		Name:            "getting-dockerfile",
		Image:           "curlimages/curl",
		ImagePullPolicy: "IfNotPresent",
		VolumeMounts: []v1.VolumeMount{
			{Name: "context", MountPath: "/context", ReadOnly: false},
		},
		Command: []string{"curl"},
		Args:    []string{"-o", "/context/Dockerfile", "-L", "https://raw.githubusercontent.com/clementreiffers/JobBuilder/main/build-worker.Dockerfile"},
	}
}

func generateCapnp() v1.Container {
	return v1.Container{
		Name:            "generating-capnp",
		Image:           "node",
		ImagePullPolicy: "IfNotPresent",
		Env:             generateAwsConfig(),
		VolumeMounts: []v1.VolumeMount{
			{Name: "s3-config", MountPath: "/root/.aws", ReadOnly: true},
			{Name: "context", MountPath: "/context"},
		},
		Command: []string{"npx", "new-capnp-generator"},
		Args: []string{
			"--bucketName=$(AWS_BUCKET)",
			"--s3Endpoint=$(AWS_ENDPOINT)",
			"--outFile=/context/config.capnp",
		},
	}
}

func generateKaniko(instance *apiv1.JobBuilder) v1.Container {
	return v1.Container{
		Name:  "kaniko",
		Image: "gcr.io/kaniko-project/executor:latest",
		Args: []string{
			"--dockerfile=Dockerfile",
			"--context=/context",
			fmt.Sprintf("--destination=%s", instance.Spec.TargetImage),
		},
		VolumeMounts: []v1.VolumeMount{
			{Name: "registry-credentials", MountPath: "/kaniko/.docker/", ReadOnly: true},
			{Name: "context", MountPath: "/context"},
		},
	}
}

func generateVolumes() []v1.Volume {
	return []v1.Volume{
		{
			Name: "registry-credentials",
			VolumeSource: v1.VolumeSource{
				Projected: &v1.ProjectedVolumeSource{
					Sources: []v1.VolumeProjection{
						{
							Secret: &v1.SecretProjection{
								LocalObjectReference: v1.LocalObjectReference{
									Name: "docker-hub",
								},
								Items: []v1.KeyToPath{
									{Key: ".dockerconfigjson", Path: "config.json"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "s3-config",
			VolumeSource: v1.VolumeSource{
				Projected: &v1.ProjectedVolumeSource{
					Sources: []v1.VolumeProjection{
						{
							Secret: &v1.SecretProjection{
								LocalObjectReference: v1.LocalObjectReference{Name: "s3-credentials"},
								Items: []v1.KeyToPath{
									{Key: "credentials", Path: "credentials"},
								},
							},
						},
						{
							ConfigMap: &v1.ConfigMapProjection{
								LocalObjectReference: v1.LocalObjectReference{Name: "aws-config"},
								Items: []v1.KeyToPath{
									{Key: "config", Path: "config"},
								},
								Optional: nil,
							},
						},
					},
				},
			},
		},
		{
			Name: "context",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}
}

func createJob(instance *apiv1.JobBuilder) batchv1.Job {
	ttl := int32(3600)
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getJobName(instance.Name),
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			//Parallelism: new(int32),
			//Completions: new(int32),
			TTLSecondsAfterFinished: &ttl,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					InitContainers: []v1.Container{
						generateDownloadFilesContainer(instance),
						generateGettingDockerfile(),
						generateCapnp(),
					},
					Containers: []v1.Container{
						generateKaniko(instance),
					},
					Volumes:       generateVolumes(),
					RestartPolicy: "Never",
				},
			},
		},
	}
}
