package postgres

import (
	"fmt"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultImage  = "postgres:16-alpine"
	PGVectorImage = "pgvector/pgvector:pg16"
	PortName      = "postgres"
	Port          = int32(5432)
	DefaultDBName = "opsmate"
)

func Deployment(config *opsmatev1alpha1.OpsMateConfig) *appsv1.Deployment {
	labels := labelsFor(config)
	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ResourceName(config),
			Namespace: config.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "postgres",
							Image: imageFor(config),
							Ports: []corev1.ContainerPort{
								{Name: PortName, ContainerPort: Port},
							},
							Env: []corev1.EnvVar{
								{Name: "POSTGRES_DB", Value: DefaultDBName},
								{Name: "POSTGRES_USER", Value: DefaultDBName},
								{
									Name: "POSTGRES_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: SecretName(config)},
											Key:                  "password",
										},
									},
								},
								{Name: "POSTGRES_SHARED_BUFFERS", Value: config.Spec.Database.SharedBuffers},
								{Name: "POSTGRES_MAX_CONNECTIONS", Value: fmt.Sprintf("%d", config.Spec.Database.MaxConnections)},
							},
						},
					},
				},
			},
		},
	}
}

func imageFor(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Database.Image != "" {
		return config.Spec.Database.Image
	}
	if config.Spec.Embedding.RequirePGVector {
		return PGVectorImage
	}
	return DefaultImage
}

func Service(config *opsmatev1alpha1.OpsMateConfig) *corev1.Service {
	labels := labelsFor(config)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ResourceName(config),
			Namespace: config.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       PortName,
					Port:       Port,
					TargetPort: intstr.FromString(PortName),
				},
			},
		},
	}
}

func ResourceName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-postgres", config.Name)
}

func SecretName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-postgres-credentials", config.Name)
}

func labelsFor(config *opsmatev1alpha1.OpsMateConfig) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "cywell-opsmate",
		"app.kubernetes.io/component":  "postgres",
		"app.kubernetes.io/instance":   config.Name,
		"app.kubernetes.io/managed-by": "cywell-opsmate-operator",
	}
}

// Package postgres will manage the PostgreSQL database used by Lightspeed-style
// conversation cache and later OpsMate operational history storage.
