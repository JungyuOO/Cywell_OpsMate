package authproxy

import (
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeploymentBuildsOAuthProxyShape(t *testing.T) {
	config := sampleConfig()

	deployment := Deployment(config)

	if deployment.Name != "sample-admin-authproxy" {
		t.Fatalf("name = %q, want sample-admin-authproxy", deployment.Name)
	}
	container := deployment.Spec.Template.Spec.Containers[0]
	if container.Image != "registry.example.com/oauth-proxy:latest" {
		t.Fatalf("image = %q, want custom image", container.Image)
	}
	assertArg(t, container.Args, "--provider=openshift")
	assertArg(t, container.Args, "--openshift-service-account=sample-admin-authproxy")
	assertArg(t, container.Args, "--upstream=https://sample-appserver:8443")
	assertArg(t, container.Args, "--pass-user-headers=true")
	assertVolume(t, deployment.Spec.Template.Spec.Volumes, "cookie-secret", "sample-cookie")
	if container.Ports[0].ContainerPort != Port {
		t.Fatalf("port = %d, want %d", container.Ports[0].ContainerPort, Port)
	}
}

func TestServiceAccountBuildsOAuthRedirectReference(t *testing.T) {
	config := sampleConfig()

	serviceAccount := ServiceAccount(config)

	if serviceAccount.Name != "sample-admin-authproxy" {
		t.Fatalf("service account name = %q", serviceAccount.Name)
	}
	annotation := serviceAccount.Annotations["serviceaccounts.openshift.io/oauth-redirectreference.sample-admin-authproxy"]
	if annotation == "" {
		t.Fatal("missing oauth redirect reference annotation")
	}
	if annotation != `{"kind":"OAuthRedirectReference","apiVersion":"v1","reference":{"kind":"Route","name":"sample-admin-authproxy"}}` {
		t.Fatalf("annotation = %q", annotation)
	}
}

func TestServiceAndRouteTargetAuthProxyOnly(t *testing.T) {
	config := sampleConfig()

	service := Service(config)
	route := Route(config)

	if service.Name != "sample-admin-authproxy" {
		t.Fatalf("service name = %q", service.Name)
	}
	if service.Spec.Selector["app.kubernetes.io/component"] != "admin-authproxy" {
		t.Fatalf("selector = %q, want admin-authproxy", service.Spec.Selector["app.kubernetes.io/component"])
	}
	if service.Annotations[ServiceCertAnnotation] != "sample-admin-authproxy-tls" {
		t.Fatalf("serving cert annotation = %q", service.Annotations[ServiceCertAnnotation])
	}
	spec := route.Object["spec"].(map[string]any)
	to := spec["to"].(map[string]any)
	if to["name"] != "sample-admin-authproxy" {
		t.Fatalf("route target = %q, want auth proxy service", to["name"])
	}
	if spec["host"] != "cyops-admin.apps.example.com" {
		t.Fatalf("host = %q", spec["host"])
	}
}

func sampleConfig() *opsmatev1alpha1.OpsMateConfig {
	return &opsmatev1alpha1.OpsMateConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample",
			Namespace: "opsmate",
		},
		Spec: opsmatev1alpha1.OpsMateConfigSpec{
			Console: opsmatev1alpha1.ConsoleSpec{
				AdminAuthProxyEnabled:         true,
				AdminAuthProxyImage:           "registry.example.com/oauth-proxy:latest",
				AdminAuthProxyCookieSecretRef: "sample-cookie",
				AdminRouteHost:                "cyops-admin.apps.example.com",
			},
		},
	}
}

func assertArg(t *testing.T, args []string, want string) {
	t.Helper()
	for _, arg := range args {
		if arg == want {
			return
		}
	}
	t.Fatalf("missing arg %q in %v", want, args)
}

func assertVolume(t *testing.T, volumes []corev1.Volume, name string, secretName string) {
	t.Helper()
	for _, volume := range volumes {
		if volume.Name != name {
			continue
		}
		if volume.Secret == nil || volume.Secret.SecretName != secretName {
			t.Fatalf("volume %s secret = %#v, want %s", name, volume.Secret, secretName)
		}
		return
	}
	t.Fatalf("missing volume %s", name)
}
