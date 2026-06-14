package controller

import (
	"context"
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestReconcileCreatesAppserverAndPostgresResources(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := opsmatev1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Lightspeed.APIBaseURL = "https://lightspeed.example.com"
	config.Spec.Lightspeed.CredentialsSecretRef = "lightspeed-secret"
	config.Spec.Database.SharedBuffers = "128MB"
	config.Spec.Database.MaxConnections = 100

	reconciler := &OpsMateConfigReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithStatusSubresource(&opsmatev1alpha1.OpsMateConfig{}).
			WithObjects(config).
			Build(),
		Scheme: scheme,
	}

	if _, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(config)}); err != nil {
		t.Fatal(err)
	}

	assertDeploymentExists(t, ctx, reconciler.Client, "sample-appserver")
	assertDeploymentExists(t, ctx, reconciler.Client, "sample-postgres")
	assertServiceExists(t, ctx, reconciler.Client, "sample-appserver")
	assertServiceExists(t, ctx, reconciler.Client, "sample-postgres")

	updated := &opsmatev1alpha1.OpsMateConfig{}
	if err := reconciler.Get(ctx, client.ObjectKeyFromObject(config), updated); err != nil {
		t.Fatal(err)
	}
	if updated.Status.OverallStatus != "Ready" {
		t.Fatalf("overall status = %q, want Ready", updated.Status.OverallStatus)
	}
	if len(updated.Status.Conditions) != 1 || updated.Status.Conditions[0].Type != "Ready" {
		t.Fatalf("conditions = %#v, want one Ready condition", updated.Status.Conditions)
	}
}

func assertDeploymentExists(t *testing.T, ctx context.Context, c client.Client, name string) {
	t.Helper()
	deployment := &appsv1.Deployment{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: "opsmate", Name: name}, deployment); err != nil {
		t.Fatalf("deployment %s missing: %v", name, err)
	}
	if len(deployment.OwnerReferences) != 1 {
		t.Fatalf("deployment %s owner references = %d, want 1", name, len(deployment.OwnerReferences))
	}
}

func assertServiceExists(t *testing.T, ctx context.Context, c client.Client, name string) {
	t.Helper()
	service := &corev1.Service{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: "opsmate", Name: name}, service); err != nil {
		t.Fatalf("service %s missing: %v", name, err)
	}
	if len(service.OwnerReferences) != 1 {
		t.Fatalf("service %s owner references = %d, want 1", name, len(service.OwnerReferences))
	}
}
