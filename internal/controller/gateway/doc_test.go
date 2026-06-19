package gateway

import (
	"strings"
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGatewayResourcesTargetNginxAndAppserver(t *testing.T) {
	config := sampleConfig()

	configMap := ConfigMap(config)
	deployment := Deployment(config)
	service := Service(config)

	if configMap.Name != "sample-gateway" {
		t.Fatalf("configmap name = %q, want sample-gateway", configMap.Name)
	}
	nginx := configMap.Data[ConfigMapKey]
	if !strings.Contains(nginx, "listen 8443 ssl") {
		t.Fatalf("nginx config missing TLS listener: %s", nginx)
	}
	if !strings.Contains(nginx, "proxy_pass https://sample-appserver.opsmate.svc:8443") {
		t.Fatalf("nginx config does not proxy appserver: %s", nginx)
	}

	container := deployment.Spec.Template.Spec.Containers[0]
	if container.Image != DefaultImage {
		t.Fatalf("image = %q, want %q", container.Image, DefaultImage)
	}
	if deployment.Spec.Template.Spec.Volumes[1].Secret.SecretName != "sample-gateway-tls" {
		t.Fatalf("tls secret = %q, want sample-gateway-tls", deployment.Spec.Template.Spec.Volumes[1].Secret.SecretName)
	}
	if service.Spec.Ports[0].Port != Port {
		t.Fatalf("service port = %d, want %d", service.Spec.Ports[0].Port, Port)
	}
	if service.Spec.Selector["app.kubernetes.io/component"] != "gateway" {
		t.Fatalf("selector component = %q, want gateway", service.Spec.Selector["app.kubernetes.io/component"])
	}
	if service.Annotations[ServiceCertAnnotation] != "sample-gateway-tls" {
		t.Fatalf("serving cert annotation = %q, want sample-gateway-tls", service.Annotations[ServiceCertAnnotation])
	}
}

func sampleConfig() *opsmatev1alpha1.OpsMateConfig {
	return &opsmatev1alpha1.OpsMateConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "sample", Namespace: "opsmate"},
	}
}
