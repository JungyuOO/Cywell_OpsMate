package console

import (
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestPluginBuildsConsolePluginShape(t *testing.T) {
	config := sampleConfig()

	plugin := Plugin(config)

	if plugin.GroupVersionKind() != ConsolePluginGroupVersionKind {
		t.Fatalf("gvk = %s, want %s", plugin.GroupVersionKind(), ConsolePluginGroupVersionKind)
	}
	if plugin.GetName() != "sample-console" {
		t.Fatalf("name = %q, want sample-console", plugin.GetName())
	}
	displayName, _, err := unstructured.NestedString(plugin.Object, "spec", "displayName")
	if err != nil {
		t.Fatal(err)
	}
	if displayName != "OpsMate" {
		t.Fatalf("displayName = %q, want OpsMate", displayName)
	}
	serviceName, _, err := unstructured.NestedString(plugin.Object, "spec", "backend", "service", "name")
	if err != nil {
		t.Fatal(err)
	}
	if serviceName != "sample-appserver" {
		t.Fatalf("service name = %q, want sample-appserver", serviceName)
	}
	serviceNamespace, _, err := unstructured.NestedString(plugin.Object, "spec", "backend", "service", "namespace")
	if err != nil {
		t.Fatal(err)
	}
	if serviceNamespace != "opsmate" {
		t.Fatalf("service namespace = %q, want opsmate", serviceNamespace)
	}
}

func TestPluginDefaultsDisplayName(t *testing.T) {
	config := sampleConfig()
	config.Spec.Console.DisplayName = ""

	plugin := Plugin(config)

	displayName, _, err := unstructured.NestedString(plugin.Object, "spec", "displayName")
	if err != nil {
		t.Fatal(err)
	}
	if displayName != DefaultDisplayName {
		t.Fatalf("displayName = %q, want %q", displayName, DefaultDisplayName)
	}
}

func TestEnabledFollowsSpec(t *testing.T) {
	config := sampleConfig()
	if !Enabled(config) {
		t.Fatal("enabled = false, want true")
	}
	config.Spec.Console.Enabled = false
	if Enabled(config) {
		t.Fatal("enabled = true, want false")
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
				Enabled:     true,
				DisplayName: "OpsMate",
			},
		},
	}
}
