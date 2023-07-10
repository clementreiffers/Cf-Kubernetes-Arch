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

func generateCmdPrebuild(scriptUrls []string, finalPath string) string {
	return fmt.Sprintf("--s3-bucket-name %s --s3-endpoint %s --s3-region fr-par --destination %s --s3-object-key %s",
		"$(AWS_BUCKET)",
		"$(AWS_ENDPOINT)",
		finalPath,
		strings.Join(scriptUrls, ","))
}

func generateDownloadFilesContainer(instance *apiv1.JobBuilder) v1.Container {
	return v1.Container{
		Name:            "download-files",
		Image:           "clementreiffers/s3-downloader-capnp-generator",
		ImagePullPolicy: "Always",
		Env:             generateAwsConfig(),
		VolumeMounts: []v1.VolumeMount{
			{Name: "s3-config", MountPath: "/root/.aws", ReadOnly: true},
			{Name: "context", MountPath: "/context"},
		},
		Command: []string{"./s3-download-files-capnp-generator"},
		Args: []string{
			"--s3-bucket-name", "$(AWS_BUCKET)",
			"--s3-endpoint", "$(AWS_ENDPOINT)",
			"--s3-region", "fr-par",
			"--destination", "/context",
			"--s3-object-key", strings.Join(instance.Spec.ScriptUrls, ","),
		},
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
		Args:    []string{"-o", "/context/Dockerfile", "-L", "https://raw.githubusercontent.com/clementreiffers/Cf-Kubernetes-Arch/main/workerd.Dockerfile"},
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
