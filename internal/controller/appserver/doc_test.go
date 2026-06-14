package appserver

import (
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeploymentBuildsAppserverShape(t *testing.T) {
	config := sampleConfig()

	deployment := Deployment(config)

	if deployment.Name != "sample-appserver" {
		t.Fatalf("name = %q, want sample-appserver", deployment.Name)
	}
	if deployment.Namespace != "opsmate" {
		t.Fatalf("namespace = %q, want opsmate", deployment.Namespace)
	}
	container := deployment.Spec.Template.Spec.Containers[0]
	if container.Image != DefaultImage {
		t.Fatalf("image = %q, want %q", container.Image, DefaultImage)
	}
	assertEnv(t, container.Env, "LIGHTSPEED_API_BASE_URL", "https://lightspeed.example.com")
	assertEnv(t, container.Env, "LIGHTSPEED_CREDENTIALS_SECRET", "lightspeed-secret")
	assertEnv(t, container.Env, "POSTGRES_SERVICE_HOST", "sample-postgres")
	assertEnv(t, container.Env, "TLS_CERT_FILE", TLSMountPath+"/tls.crt")
	assertEnv(t, container.Env, "TLS_KEY_FILE", TLSMountPath+"/tls.key")
	if deployment.Spec.Template.Spec.Volumes[0].Secret.SecretName != "sample-appserver-tls" {
		t.Fatalf("tls secret = %q, want sample-appserver-tls", deployment.Spec.Template.Spec.Volumes[0].Secret.SecretName)
	}
	if container.VolumeMounts[0].MountPath != TLSMountPath {
		t.Fatalf("tls mount path = %q, want %q", container.VolumeMounts[0].MountPath, TLSMountPath)
	}
	if container.Ports[0].ContainerPort != Port {
		t.Fatalf("container port = %d, want %d", container.Ports[0].ContainerPort, Port)
	}
}

func TestServiceTargetsAppserverPort(t *testing.T) {
	config := sampleConfig()

	service := Service(config)

	if service.Name != "sample-appserver" {
		t.Fatalf("name = %q, want sample-appserver", service.Name)
	}
	if service.Spec.Ports[0].Port != Port {
		t.Fatalf("service port = %d, want %d", service.Spec.Ports[0].Port, Port)
	}
	if service.Spec.Ports[0].TargetPort.StrVal != PortName {
		t.Fatalf("target port = %q, want %q", service.Spec.Ports[0].TargetPort.StrVal, PortName)
	}
	if service.Spec.Selector["app.kubernetes.io/component"] != "appserver" {
		t.Fatalf("component selector = %q, want appserver", service.Spec.Selector["app.kubernetes.io/component"])
	}
	if service.Annotations[ServiceCertAnnotation] != "sample-appserver-tls" {
		t.Fatalf("serving cert annotation = %q, want sample-appserver-tls", service.Annotations[ServiceCertAnnotation])
	}
}

func sampleConfig() *opsmatev1alpha1.OpsMateConfig {
	return &opsmatev1alpha1.OpsMateConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample",
			Namespace: "opsmate",
		},
		Spec: opsmatev1alpha1.OpsMateConfigSpec{
			Lightspeed: opsmatev1alpha1.LightspeedSpec{
				APIBaseURL:           "https://lightspeed.example.com",
				CredentialsSecretRef: "lightspeed-secret",
				DefaultProvider:      "openai",
				DefaultModel:         "gpt-4.1",
			},
		},
	}
}

func assertEnv(t *testing.T, env []corev1.EnvVar, name string, want string) {
	t.Helper()
	for _, item := range env {
		if item.Name == name {
			if item.Value != want {
				t.Fatalf("%s = %q, want %q", name, item.Value, want)
			}
			return
		}
	}
	t.Fatalf("missing env %s", name)
}
