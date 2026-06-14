package appserver

import (
	"fmt"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultImage = "ghcr.io/jungyuoo/cywell-opsmate-appserver:latest"
	PortName     = "http"
	Port         = int32(8080)
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
							Name:  "appserver",
							Image: DefaultImage,
							Ports: []corev1.ContainerPort{
								{Name: PortName, ContainerPort: Port},
							},
							Env: []corev1.EnvVar{
								{Name: "LIGHTSPEED_API_BASE_URL", Value: config.Spec.Lightspeed.APIBaseURL},
								{Name: "LIGHTSPEED_CREDENTIALS_SECRET", Value: config.Spec.Lightspeed.CredentialsSecretRef},
								{Name: "LIGHTSPEED_DEFAULT_PROVIDER", Value: config.Spec.Lightspeed.DefaultProvider},
								{Name: "LIGHTSPEED_DEFAULT_MODEL", Value: config.Spec.Lightspeed.DefaultModel},
								{Name: "POSTGRES_SERVICE_HOST", Value: fmt.Sprintf("%s-postgres", config.Name)},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: ptr(false),
								RunAsNonRoot:             ptr(true),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: ptr(true),
					},
				},
			},
		},
	}
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
	return fmt.Sprintf("%s-appserver", config.Name)
}

func labelsFor(config *opsmatev1alpha1.OpsMateConfig) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "cywell-opsmate",
		"app.kubernetes.io/component":  "appserver",
		"app.kubernetes.io/instance":   config.Name,
		"app.kubernetes.io/managed-by": "cywell-opsmate-operator",
	}
}

func ptr[T any](value T) *T {
	return &value
}

// Package appserver will manage the OpsMate backend service that proxies
// Lightspeed API calls and hosts AIOps/RAG integration points.
