package console

import (
	"fmt"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/appserver"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	DefaultDisplayName = "Cywell OpsMate"
)

var ConsolePluginGroupVersionKind = schema.GroupVersionKind{
	Group:   "console.openshift.io",
	Version: "v1",
	Kind:    "ConsolePlugin",
}

func Enabled(config *opsmatev1alpha1.OpsMateConfig) bool {
	return config.Spec.Console.Enabled
}

func Plugin(config *opsmatev1alpha1.OpsMateConfig) *unstructured.Unstructured {
	displayName := config.Spec.Console.DisplayName
	if displayName == "" {
		displayName = DefaultDisplayName
	}

	object := &unstructured.Unstructured{}
	object.SetGroupVersionKind(ConsolePluginGroupVersionKind)
	object.SetName(ResourceName(config))
	object.SetLabels(labelsFor(config))
	object.Object["spec"] = map[string]any{
		"displayName": displayName,
		"backend": map[string]any{
			"type": "Service",
			"service": map[string]any{
				"name":      appserver.ResourceName(config),
				"namespace": config.Namespace,
				"port":      int64(appserver.Port),
				"basePath":  "/",
			},
		},
	}
	return object
}

func ResourceName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-console", config.Name)
}

func labelsFor(config *opsmatev1alpha1.OpsMateConfig) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "cywell-opsmate",
		"app.kubernetes.io/component":  "console-plugin",
		"app.kubernetes.io/instance":   config.Name,
		"app.kubernetes.io/managed-by": "cywell-opsmate-operator",
	}
}

// Package console will manage the OpenShift ConsolePlugin resources and the
// console plugin deployment served inside OpenShift Web Console.
