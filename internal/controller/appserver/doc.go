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
	DefaultImage          = "ghcr.io/jungyuoo/cywell-opsmate-appserver:latest"
	PortName              = "https"
	Port                  = int32(8443)
	ServiceCertAnnotation = "service.beta.openshift.io/serving-cert-secret-name"
	TLSMountPath          = "/var/run/secrets/cywell-opsmate/tls"
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
							Env: appserverEnv(config),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "serving-cert",
									MountPath: TLSMountPath,
									ReadOnly:  true,
								},
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
					Volumes: []corev1.Volume{
						{
							Name: "serving-cert",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: TLSSecretName(config),
								},
							},
						},
					},
				},
			},
		},
	}
}

func appserverEnv(config *opsmatev1alpha1.OpsMateConfig) []corev1.EnvVar {
	env := []corev1.EnvVar{
		{Name: "LIGHTSPEED_API_BASE_URL", Value: config.Spec.Lightspeed.APIBaseURL},
		{Name: "LIGHTSPEED_CREDENTIALS_SECRET", Value: config.Spec.Lightspeed.CredentialsSecretRef},
		{Name: "LIGHTSPEED_DEFAULT_PROVIDER", Value: config.Spec.Lightspeed.DefaultProvider},
		{Name: "LIGHTSPEED_DEFAULT_MODEL", Value: config.Spec.Lightspeed.DefaultModel},
		{Name: "CYOPS_EMBEDDING_ENDPOINT", Value: config.Spec.Embedding.EndpointURL},
		{Name: "CYOPS_EMBEDDING_MODEL", Value: config.Spec.Embedding.Model},
		{Name: "CYOPS_EMBEDDING_DIMENSIONS", Value: embeddingDimensions(config)},
		{Name: "CYOPS_PGVECTOR_REQUIRED", Value: embeddingPGVectorRequired(config)},
		{Name: "CYOPS_RETRIEVAL_MODE", Value: embeddingRetrievalMode(config)},
		{Name: "CYOPS_RETRIEVAL_SLOW_THRESHOLD_MS", Value: embeddingRetrievalSlowMillis(config)},
		{Name: "POSTGRES_SERVICE_HOST", Value: fmt.Sprintf("%s-postgres", config.Name)},
		{Name: "TLS_CERT_FILE", Value: TLSMountPath + "/tls.crt"},
		{Name: "TLS_KEY_FILE", Value: TLSMountPath + "/tls.key"},
	}
	if config.Spec.Embedding.CredentialsSecretRef != "" {
		key := config.Spec.Embedding.CredentialsSecretKey
		if key == "" {
			key = "token"
		}
		env = append(env, corev1.EnvVar{
			Name: "CYOPS_EMBEDDING_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: config.Spec.Embedding.CredentialsSecretRef},
					Key:                  key,
				},
			},
		})
	}
	if config.Spec.Database.DSNSecretRef != "" {
		key := config.Spec.Database.DSNSecretKey
		if key == "" {
			key = "dsn"
		}
		env = append(env, corev1.EnvVar{
			Name: "CYOPS_POSTGRES_DSN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: config.Spec.Database.DSNSecretRef},
					Key:                  key,
				},
			},
		})
	}
	return env
}

func embeddingDimensions(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Embedding.Dimensions <= 0 {
		return ""
	}
	return fmt.Sprintf("%d", config.Spec.Embedding.Dimensions)
}

func embeddingPGVectorRequired(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Embedding.RequirePGVector {
		return "true"
	}
	return "false"
}

func embeddingRetrievalMode(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Embedding.RetrievalMode == "" {
		return "bytea"
	}
	return config.Spec.Embedding.RetrievalMode
}

func embeddingRetrievalSlowMillis(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Embedding.RetrievalSlowMillis <= 0 {
		return ""
	}
	return fmt.Sprintf("%d", config.Spec.Embedding.RetrievalSlowMillis)
}

func Service(config *opsmatev1alpha1.OpsMateConfig) *corev1.Service {
	labels := labelsFor(config)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ResourceName(config),
			Namespace:   config.Namespace,
			Labels:      labels,
			Annotations: map[string]string{ServiceCertAnnotation: TLSSecretName(config)},
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

func TLSSecretName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-appserver-tls", config.Name)
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
