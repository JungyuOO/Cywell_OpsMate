package controller

import (
	"context"
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestReconcileCreatesAppserverAndPostgresResources(t *testing.T) {
	ctx := context.Background()
	scheme := testScheme(t)

	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Lightspeed.APIBaseURL = "https://lightspeed.example.com"
	config.Spec.Lightspeed.CredentialsSecretRef = "lightspeed-secret"
	config.Spec.Database.SharedBuffers = "128MB"
	config.Spec.Database.MaxConnections = 100
	config.Spec.Console.Enabled = true
	config.Spec.Console.DisplayName = "OpsMate"
	config.Spec.Console.AdminAuthProxyEnabled = true
	config.Spec.Console.AdminAuthProxyCookieSecretRef = "sample-cookie"

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
	assertDeploymentExists(t, ctx, reconciler.Client, "sample-admin-authproxy")
	assertServiceExists(t, ctx, reconciler.Client, "sample-appserver")
	assertServiceExists(t, ctx, reconciler.Client, "sample-postgres")
	assertServiceExists(t, ctx, reconciler.Client, "sample-admin-authproxy")
	assertServiceAccountExists(t, ctx, reconciler.Client, "sample-admin-authproxy")
	assertConsolePluginExists(t, ctx, reconciler.Client, "sample-console")
	assertRouteExists(t, ctx, reconciler.Client, "sample-admin-authproxy")

	updated := &opsmatev1alpha1.OpsMateConfig{}
	if err := reconciler.Get(ctx, client.ObjectKeyFromObject(config), updated); err != nil {
		t.Fatal(err)
	}
	if updated.Status.OverallStatus != "Ready" {
		t.Fatalf("overall status = %q, want Ready", updated.Status.OverallStatus)
	}
	assertCondition(t, updated.Status.Conditions, "Ready", "True", "ResourcesReconciled")
	assertCondition(t, updated.Status.Conditions, "PostgresDSNConfigured", "False", "SecretReferenceMissing")
	assertCondition(t, updated.Status.Conditions, "PGVectorRequired", "False", "NotRequiredBySpec")
	assertCondition(t, updated.Status.Conditions, "PGVectorMigrationApproved", "False", "ApprovalRequired")
	assertCondition(t, updated.Status.Conditions, "PGVectorReady", "True", "NotRequired")
	assertCondition(t, updated.Status.Conditions, "ReembeddingReady", "True", "ReembeddingIdle")
	assertCondition(t, updated.Status.Conditions, "RetrievalModeReady", "True", "BYTEAFallback")
}

func TestReconcileAppliesCompletedPGVectorMigrationJobStatus(t *testing.T) {
	ctx := context.Background()
	scheme := testScheme(t)

	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Database.PGVectorMigrationApproved = true
	config.Spec.Database.DSNSecretRef = "postgres-dsn"
	config.Spec.Embedding.RequirePGVector = true
	config.Spec.Embedding.RetrievalMode = "pgvector"
	config.Spec.Embedding.Dimensions = 768

	job := &batchv1.Job{}
	job.Name = "sample-pgvector-migration"
	job.Namespace = "opsmate"
	job.Status.Conditions = []batchv1.JobCondition{
		{Type: batchv1.JobComplete, Status: corev1.ConditionTrue},
	}

	reconciler := &OpsMateConfigReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithStatusSubresource(&opsmatev1alpha1.OpsMateConfig{}).
			WithObjects(config, job).
			Build(),
		Scheme: scheme,
	}

	if _, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(config)}); err != nil {
		t.Fatal(err)
	}

	updated := &opsmatev1alpha1.OpsMateConfig{}
	if err := reconciler.Get(ctx, client.ObjectKeyFromObject(config), updated); err != nil {
		t.Fatal(err)
	}
	if !updated.Status.PGVectorReady {
		t.Fatal("pgVectorReady = false, want true")
	}
	assertCondition(t, updated.Status.Conditions, "PGVectorReady", "True", "RuntimeCheckPassed")
}

func TestStatusConditionsExposePGVectorConfiguration(t *testing.T) {
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Database.DSNSecretRef = "postgres-dsn"
	config.Spec.Database.PGVectorMigrationApproved = true
	config.Spec.Embedding.RequirePGVector = true
	config.Spec.Embedding.RetrievalMode = "pgvector"
	config.Status.PGVectorReady = true

	conditions := statusConditions(config)

	if got := overallStatus(conditions); got != "Ready" {
		t.Fatalf("overall status = %q, want Ready", got)
	}
	assertCondition(t, conditions, "PostgresDSNConfigured", "True", "SecretReferenceConfigured")
	assertCondition(t, conditions, "PGVectorRequired", "True", "RequiredBySpec")
	assertCondition(t, conditions, "PGVectorMigrationApproved", "True", "ApprovedBySpec")
	assertCondition(t, conditions, "PGVectorReady", "True", "RuntimeCheckPassed")
	assertCondition(t, conditions, "RetrievalModeReady", "True", "PGVectorMode")
}

func TestStatusConditionsDegradeFailedReembedding(t *testing.T) {
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Status.Reembedding.Failed = 1

	conditions := statusConditions(config)

	if got := overallStatus(conditions); got != "Degraded" {
		t.Fatalf("overall status = %q, want Degraded", got)
	}
	assertCondition(t, conditions, "ReembeddingReady", "False", "ReembeddingFailed")
}

func TestStatusConditionsDegradeFailedPGVectorEvidence(t *testing.T) {
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Spec.Database.DSNSecretRef = "postgres-dsn"
	config.Spec.Embedding.RequirePGVector = true
	config.Status.PGVectorLastError = "pgvector migration job failed"

	conditions := statusConditions(config)

	if got := overallStatus(conditions); got != "Degraded" {
		t.Fatalf("overall status = %q, want Degraded", got)
	}
	assertCondition(t, conditions, "PGVectorReady", "False", "RuntimeCheckFailed")
}

func TestStatusConditionsDegradeInvalidPGVectorMode(t *testing.T) {
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Spec.Embedding.RetrievalMode = "pgvector"

	conditions := statusConditions(config)

	if got := overallStatus(conditions); got != "Degraded" {
		t.Fatalf("overall status = %q, want Degraded", got)
	}
	assertCondition(t, conditions, "RetrievalModeReady", "False", "PGVectorNotRequired")
}

func TestReconcileSkipsDisabledConsolePlugin(t *testing.T) {
	ctx := context.Background()
	scheme := testScheme(t)
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Console.Enabled = false

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

	plugin := consolePluginObject()
	if err := reconciler.Get(ctx, client.ObjectKey{Name: "sample-console"}, plugin); err == nil {
		t.Fatal("ConsolePlugin exists, want skipped")
	}
}

func TestReconcileUpdatesExistingConsolePlugin(t *testing.T) {
	ctx := context.Background()
	scheme := testScheme(t)
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Console.Enabled = true
	config.Spec.Console.DisplayName = "OpsMate"

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
	updatedConfig := &opsmatev1alpha1.OpsMateConfig{}
	if err := reconciler.Get(ctx, client.ObjectKeyFromObject(config), updatedConfig); err != nil {
		t.Fatal(err)
	}
	updatedConfig.Spec.Console.DisplayName = "CYOps"
	if err := reconciler.Update(ctx, updatedConfig); err != nil {
		t.Fatal(err)
	}
	if _, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(config)}); err != nil {
		t.Fatal(err)
	}

	plugin := consolePluginObject()
	if err := reconciler.Get(ctx, client.ObjectKey{Name: "sample-console"}, plugin); err != nil {
		t.Fatal(err)
	}
	displayName, _, err := unstructured.NestedString(plugin.Object, "spec", "displayName")
	if err != nil {
		t.Fatal(err)
	}
	if displayName != "CYOps" {
		t.Fatalf("displayName = %q, want CYOps", displayName)
	}
}

func TestReconcilePreservesServiceAnnotations(t *testing.T) {
	ctx := context.Background()
	scheme := testScheme(t)
	config := &opsmatev1alpha1.OpsMateConfig{}
	config.Name = "sample"
	config.Namespace = "opsmate"
	config.Spec.Lightspeed.APIBaseURL = "https://lightspeed.example.com"
	config.Spec.Console.Enabled = true

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

	service := &corev1.Service{}
	if err := reconciler.Get(ctx, client.ObjectKey{Namespace: "opsmate", Name: "sample-appserver"}, service); err != nil {
		t.Fatal(err)
	}
	if service.Annotations["service.beta.openshift.io/serving-cert-secret-name"] != "sample-appserver-tls" {
		t.Fatalf("serving cert annotation = %q", service.Annotations["service.beta.openshift.io/serving-cert-secret-name"])
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

func assertServiceAccountExists(t *testing.T, ctx context.Context, c client.Client, name string) {
	t.Helper()
	serviceAccount := &corev1.ServiceAccount{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: "opsmate", Name: name}, serviceAccount); err != nil {
		t.Fatalf("service account %s missing: %v", name, err)
	}
	if serviceAccount.Annotations["serviceaccounts.openshift.io/oauth-redirectreference."+name] == "" {
		t.Fatalf("service account %s missing oauth redirect annotation", name)
	}
}

func assertConsolePluginExists(t *testing.T, ctx context.Context, c client.Client, name string) {
	t.Helper()
	plugin := consolePluginObject()
	if err := c.Get(ctx, client.ObjectKey{Name: name}, plugin); err != nil {
		t.Fatalf("console plugin %s missing: %v", name, err)
	}
}

func assertRouteExists(t *testing.T, ctx context.Context, c client.Client, name string) {
	t.Helper()
	route := routeObject()
	if err := c.Get(ctx, client.ObjectKey{Namespace: "opsmate", Name: name}, route); err != nil {
		t.Fatalf("route %s missing: %v", name, err)
	}
}

func assertCondition(t *testing.T, conditions []metav1.Condition, conditionType string, status string, reason string) {
	t.Helper()
	for _, condition := range conditions {
		if condition.Type != conditionType {
			continue
		}
		if string(condition.Status) != status {
			t.Fatalf("%s status = %s, want %s", conditionType, condition.Status, status)
		}
		if condition.Reason != reason {
			t.Fatalf("%s reason = %q, want %q", conditionType, condition.Reason, reason)
		}
		return
	}
	t.Fatalf("missing condition %s in %#v", conditionType, conditions)
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := opsmatev1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "console.openshift.io",
		Version: "v1",
		Kind:    "ConsolePlugin",
	}, consolePluginObject())
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{
		Group:   "route.openshift.io",
		Version: "v1",
		Kind:    "Route",
	}, routeObject())
	return scheme
}

func consolePluginObject() *unstructured.Unstructured {
	plugin := &unstructured.Unstructured{}
	plugin.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "console.openshift.io",
		Version: "v1",
		Kind:    "ConsolePlugin",
	})
	return plugin
}

func routeObject() *unstructured.Unstructured {
	route := &unstructured.Unstructured{}
	route.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "route.openshift.io",
		Version: "v1",
		Kind:    "Route",
	})
	return route
}
