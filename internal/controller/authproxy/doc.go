package authproxy

import (
	"fmt"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/appserver"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultImage          = "quay.io/openshift/origin-oauth-proxy:latest"
	PortName              = "https"
	Port                  = int32(8443)
	ServiceCertAnnotation = "service.beta.openshift.io/serving-cert-secret-name"
	TLSMountPath          = "/var/run/secrets/cywell-opsmate/authproxy-tls"
	CookieMountPath       = "/var/run/secrets/cywell-opsmate/oauth-cookie"
	CookieSecretKey       = "session_secret"
)

func Enabled(config *opsmatev1alpha1.OpsMateConfig) bool {
	return config.Spec.Console.AdminAuthProxyEnabled
}

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
					ServiceAccountName: ServiceAccountName(config),
					Containers: []corev1.Container{{
						Name:  "oauth-proxy",
						Image: imageFor(config),
						Args: []string{
							"--provider=openshift",
							"--openshift-service-account=" + ServiceAccountName(config),
							"--https-address=:8443",
							"--upstream=https://" + appserver.ResourceName(config) + ":8443",
							"--tls-cert=" + TLSMountPath + "/tls.crt",
							"--tls-key=" + TLSMountPath + "/tls.key",
							"--cookie-secret-file=" + CookieMountPath + "/" + CookieSecretKey,
							"--pass-user-headers=true",
							"--pass-access-token=false",
						},
						Ports: []corev1.ContainerPort{{Name: PortName, ContainerPort: Port}},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "serving-cert", MountPath: TLSMountPath, ReadOnly: true},
							{Name: "cookie-secret", MountPath: CookieMountPath, ReadOnly: true},
						},
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: ptr(false),
							RunAsNonRoot:             ptr(true),
							Capabilities:             &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
						},
					}},
					SecurityContext: &corev1.PodSecurityContext{RunAsNonRoot: ptr(true)},
					Volumes: []corev1.Volume{
						{
							Name: "serving-cert",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{SecretName: TLSSecretName(config)},
							},
						},
						{
							Name: "cookie-secret",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{SecretName: cookieSecretName(config)},
							},
						},
					},
				},
			},
		},
	}
}

func ServiceAccount(config *opsmatev1alpha1.OpsMateConfig) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceAccountName(config),
			Namespace: config.Namespace,
			Labels:    labelsFor(config),
			Annotations: map[string]string{
				"serviceaccounts.openshift.io/oauth-redirectreference." + ResourceName(config): oauthRedirectReference(config),
			},
		},
	}
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
			Ports: []corev1.ServicePort{{
				Name:       PortName,
				Port:       Port,
				TargetPort: intstr.FromString(PortName),
			}},
		},
	}
}

func Route(config *opsmatev1alpha1.OpsMateConfig) *unstructured.Unstructured {
	route := &unstructured.Unstructured{}
	route.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "route.openshift.io",
		Version: "v1",
		Kind:    "Route",
	})
	route.SetName(ResourceName(config))
	route.SetNamespace(config.Namespace)
	route.SetLabels(labelsFor(config))

	spec := map[string]any{
		"to": map[string]any{
			"kind": "Service",
			"name": ResourceName(config),
		},
		"port": map[string]any{"targetPort": PortName},
		"tls": map[string]any{
			"termination":                   "reencrypt",
			"insecureEdgeTerminationPolicy": "Redirect",
		},
	}
	if config.Spec.Console.AdminRouteHost != "" {
		spec["host"] = config.Spec.Console.AdminRouteHost
	}
	route.Object["spec"] = spec
	return route
}

func ResourceName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-admin-authproxy", config.Name)
}

func ServiceAccountName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-admin-authproxy", config.Name)
}

func TLSSecretName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-admin-authproxy-tls", config.Name)
}

func oauthRedirectReference(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf(`{"kind":"OAuthRedirectReference","apiVersion":"v1","reference":{"kind":"Route","name":"%s"}}`, ResourceName(config))
}

func imageFor(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Console.AdminAuthProxyImage != "" {
		return config.Spec.Console.AdminAuthProxyImage
	}
	return DefaultImage
}

func cookieSecretName(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Console.AdminAuthProxyCookieSecretRef != "" {
		return config.Spec.Console.AdminAuthProxyCookieSecretRef
	}
	return fmt.Sprintf("%s-admin-oauth-cookie", config.Name)
}

func labelsFor(config *opsmatev1alpha1.OpsMateConfig) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "cywell-opsmate",
		"app.kubernetes.io/component":  "admin-authproxy",
		"app.kubernetes.io/instance":   config.Name,
		"app.kubernetes.io/managed-by": "cywell-opsmate-operator",
	}
}

func ptr[T any](value T) *T {
	return &value
}

// Package authproxy builds the OpenShift OAuth-protected admin entry point.
