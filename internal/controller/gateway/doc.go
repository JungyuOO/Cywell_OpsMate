package gateway

import (
	"fmt"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/appserver"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultImage          = "nginxinc/nginx-unprivileged:1.27-alpine"
	PortName              = "https"
	Port                  = int32(8443)
	ConfigMapKey          = "nginx.conf"
	TLSMountPath          = "/var/run/secrets/cywell-opsmate/gateway-tls"
	ConfigMountPath       = "/etc/nginx/nginx.conf"
	ServiceCertAnnotation = appserver.ServiceCertAnnotation
)

func ConfigMap(config *opsmatev1alpha1.OpsMateConfig) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ResourceName(config),
			Namespace: config.Namespace,
			Labels:    labelsFor(config),
		},
		Data: map[string]string{
			ConfigMapKey: nginxConfig(config),
		},
	}
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
					SecurityContext: &corev1.PodSecurityContext{RunAsNonRoot: ptr(true)},
					Containers: []corev1.Container{
						{
							Name:  "gateway",
							Image: DefaultImage,
							Ports: []corev1.ContainerPort{
								{Name: PortName, ContainerPort: Port},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "nginx-config",
									MountPath: ConfigMountPath,
									SubPath:   ConfigMapKey,
									ReadOnly:  true,
								},
								{
									Name:      "serving-cert",
									MountPath: TLSMountPath,
									ReadOnly:  true,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: ptr(false),
								RunAsNonRoot:             ptr(true),
								Capabilities:             &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "nginx-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: ResourceName(config)},
								},
							},
						},
						{
							Name: "serving-cert",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{SecretName: TLSSecretName(config)},
							},
						},
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
			Name:        ResourceName(config),
			Namespace:   config.Namespace,
			Labels:      labels,
			Annotations: map[string]string{ServiceCertAnnotation: TLSSecretName(config)},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{Name: PortName, Port: Port, TargetPort: intstr.FromString(PortName)},
			},
		},
	}
}

func ResourceName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-gateway", config.Name)
}

func TLSSecretName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-gateway-tls", config.Name)
}

func nginxConfig(config *opsmatev1alpha1.OpsMateConfig) string {
	upstream := fmt.Sprintf("https://%s.%s.svc:8443", appserver.ResourceName(config), config.Namespace)
	return fmt.Sprintf(`pid /tmp/nginx.pid;
events {}
http {
  access_log /dev/stdout;
  error_log /dev/stderr warn;
  client_body_temp_path /tmp/client_body;
  proxy_temp_path /tmp/proxy;
  fastcgi_temp_path /tmp/fastcgi;
  uwsgi_temp_path /tmp/uwsgi;
  scgi_temp_path /tmp/scgi;

  server {
    listen 8443 ssl;
    ssl_certificate %s/tls.crt;
    ssl_certificate_key %s/tls.key;

    location / {
      proxy_pass %s;
      proxy_ssl_verify off;
      proxy_set_header Host $host;
      proxy_set_header X-Forwarded-Proto https;
      proxy_set_header X-Forwarded-Host $host;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
  }
}
`, TLSMountPath, TLSMountPath, upstream)
}

func labelsFor(config *opsmatev1alpha1.OpsMateConfig) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "cywell-opsmate",
		"app.kubernetes.io/component":  "gateway",
		"app.kubernetes.io/instance":   config.Name,
		"app.kubernetes.io/managed-by": "cywell-opsmate-operator",
	}
}

func ptr[T any](value T) *T {
	return &value
}

// Package gateway manages the nginx front door used by OpenShift ConsolePlugin.
