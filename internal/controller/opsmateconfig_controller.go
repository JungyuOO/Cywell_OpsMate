package controller

import (
	"context"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/appserver"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/authproxy"
	consoleplugin "github.com/JungyuOO/Cywell_OpsMate/internal/controller/console"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/postgres"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type OpsMateConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func SetupOpsMateConfigReconciler(mgr ctrl.Manager) error {
	return (&OpsMateConfigReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr)
}

// +kubebuilder:rbac:groups=opsmate.cywell.io,resources=opsmateconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=opsmate.cywell.io,resources=opsmateconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=console.openshift.io,resources=consoleplugins,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch
func (r *OpsMateConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	config := &opsmatev1alpha1.OpsMateConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, object := range []client.Object{
		appserver.Deployment(config),
		appserver.Service(config),
		postgres.Deployment(config),
		postgres.Service(config),
	} {
		if err := controllerutil.SetControllerReference(config, object, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.reconcileObject(ctx, object); err != nil {
			return ctrl.Result{}, err
		}
	}

	if consoleplugin.Enabled(config) {
		plugin := consoleplugin.Plugin(config)
		if err := r.reconcileObject(ctx, plugin); err != nil {
			return ctrl.Result{}, err
		}
	}
	if authproxy.Enabled(config) {
		for _, object := range []client.Object{
			authproxy.ServiceAccount(config),
			authproxy.Deployment(config),
			authproxy.Service(config),
			authproxy.Route(config),
		} {
			if err := controllerutil.SetControllerReference(config, object, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.reconcileObject(ctx, object); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	if config.Spec.Database.PGVectorMigrationApproved {
		job, err := postgres.PGVectorMigrationJob(config)
		if err == nil {
			if err := controllerutil.SetControllerReference(config, job, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.reconcileObject(ctx, job); err != nil {
				return ctrl.Result{}, err
			}
			current := &batchv1.Job{}
			if err := r.Get(ctx, client.ObjectKeyFromObject(job), current); err != nil {
				return ctrl.Result{}, err
			}
			postgres.ApplyPGVectorMigrationJobStatus(config, current)
		} else {
			config.Status.PGVectorLastError = err.Error()
		}
	}

	config.Status.Conditions = statusConditions(config)
	config.Status.OverallStatus = overallStatus(config.Status.Conditions)
	if err := r.Status().Update(ctx, config); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func statusConditions(config *opsmatev1alpha1.OpsMateConfig) []metav1.Condition {
	conditions := []metav1.Condition{}
	ready := condition(
		"Ready",
		metav1.ConditionTrue,
		config.Generation,
		"ResourcesReconciled",
		"Appserver and PostgreSQL resources are reconciled.",
	)
	conditions = upsertCondition(conditions, ready)
	conditions = upsertCondition(conditions, postgresDSNConfiguredCondition(config))
	conditions = upsertCondition(conditions, pgVectorRequiredCondition(config))
	conditions = upsertCondition(conditions, pgVectorMigrationApprovedCondition(config))
	conditions = upsertCondition(conditions, pgVectorReadyCondition(config))
	conditions = upsertCondition(conditions, reembeddingCondition(config))
	conditions = upsertCondition(conditions, retrievalModeReadyCondition(config))
	return conditions
}

func condition(conditionType string, status metav1.ConditionStatus, generation int64, reason string, message string) metav1.Condition {
	return metav1.Condition{
		Type:               conditionType,
		Status:             status,
		ObservedGeneration: generation,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func postgresDSNConfiguredCondition(config *opsmatev1alpha1.OpsMateConfig) metav1.Condition {
	if config.Spec.Database.DSNSecretRef != "" {
		return condition(
			"PostgresDSNConfigured",
			metav1.ConditionTrue,
			config.Generation,
			"SecretReferenceConfigured",
			"PostgreSQL DSN is configured through a Secret reference.",
		)
	}
	return condition(
		"PostgresDSNConfigured",
		metav1.ConditionFalse,
		config.Generation,
		"SecretReferenceMissing",
		"PostgreSQL DSN Secret reference is not configured.",
	)
}

func pgVectorRequiredCondition(config *opsmatev1alpha1.OpsMateConfig) metav1.Condition {
	if config.Spec.Embedding.RequirePGVector {
		return condition(
			"PGVectorRequired",
			metav1.ConditionTrue,
			config.Generation,
			"RequiredBySpec",
			"pgvector startup validation is required by spec.",
		)
	}
	return condition(
		"PGVectorRequired",
		metav1.ConditionFalse,
		config.Generation,
		"NotRequiredBySpec",
		"pgvector startup validation is not required by spec.",
	)
}

func pgVectorMigrationApprovedCondition(config *opsmatev1alpha1.OpsMateConfig) metav1.Condition {
	if config.Spec.Database.PGVectorMigrationApproved {
		return condition(
			"PGVectorMigrationApproved",
			metav1.ConditionTrue,
			config.Generation,
			"ApprovedBySpec",
			"pgvector migration is explicitly approved by spec.",
		)
	}
	return condition(
		"PGVectorMigrationApproved",
		metav1.ConditionFalse,
		config.Generation,
		"ApprovalRequired",
		"pgvector migration is not approved; reconciliation will not apply it automatically.",
	)
}

func pgVectorReadyCondition(config *opsmatev1alpha1.OpsMateConfig) metav1.Condition {
	if !config.Spec.Embedding.RequirePGVector {
		return condition(
			"PGVectorReady",
			metav1.ConditionTrue,
			config.Generation,
			"NotRequired",
			"pgvector is not required for the current retrieval mode.",
		)
	}
	if config.Spec.Database.DSNSecretRef == "" {
		return condition(
			"PGVectorReady",
			metav1.ConditionFalse,
			config.Generation,
			"DSNSecretReferenceMissing",
			"pgvector readiness cannot be checked until a PostgreSQL DSN Secret reference is configured.",
		)
	}
	if config.Status.PGVectorReady {
		return condition(
			"PGVectorReady",
			metav1.ConditionTrue,
			config.Generation,
			"RuntimeCheckPassed",
			"pgvector runtime validation has passed.",
		)
	}
	if config.Status.PGVectorLastError != "" {
		return condition(
			"PGVectorReady",
			metav1.ConditionFalse,
			config.Generation,
			"RuntimeCheckFailed",
			config.Status.PGVectorLastError,
		)
	}
	return condition(
		"PGVectorReady",
		metav1.ConditionUnknown,
		config.Generation,
		"RuntimeCheckPending",
		"pgvector readiness is validated by appserver startup and live smoke tests.",
	)
}

func reembeddingCondition(config *opsmatev1alpha1.OpsMateConfig) metav1.Condition {
	if config.Status.Reembedding.Running {
		return condition(
			"ReembeddingReady",
			metav1.ConditionUnknown,
			config.Generation,
			"ReembeddingRunning",
			"Document re-embedding is running.",
		)
	}
	if config.Status.Reembedding.Failed > 0 {
		return condition(
			"ReembeddingReady",
			metav1.ConditionFalse,
			config.Generation,
			"ReembeddingFailed",
			"Document re-embedding completed with failures.",
		)
	}
	return condition(
		"ReembeddingReady",
		metav1.ConditionTrue,
		config.Generation,
		"ReembeddingIdle",
		"Document re-embedding is not running.",
	)
}

func retrievalModeReadyCondition(config *opsmatev1alpha1.OpsMateConfig) metav1.Condition {
	switch config.Spec.Embedding.RetrievalMode {
	case "", "bytea":
		return condition(
			"RetrievalModeReady",
			metav1.ConditionTrue,
			config.Generation,
			"BYTEAFallback",
			"Retrieval mode is configured for BYTEA fallback.",
		)
	case "pgvector":
		if config.Spec.Embedding.RequirePGVector {
			return condition(
				"RetrievalModeReady",
				metav1.ConditionTrue,
				config.Generation,
				"PGVectorMode",
				"Retrieval mode is configured for pgvector with pgvector required.",
			)
		}
		return condition(
			"RetrievalModeReady",
			metav1.ConditionFalse,
			config.Generation,
			"PGVectorNotRequired",
			"retrievalMode=pgvector requires requirePGVector=true.",
		)
	default:
		return condition(
			"RetrievalModeReady",
			metav1.ConditionFalse,
			config.Generation,
			"UnsupportedRetrievalMode",
			"retrievalMode must be bytea or pgvector.",
		)
	}
}

func overallStatus(conditions []metav1.Condition) string {
	for _, item := range conditions {
		if item.Status == metav1.ConditionFalse && (item.Type == "PGVectorReady" || item.Type == "RetrievalModeReady" || item.Type == "ReembeddingReady") {
			return "Degraded"
		}
	}
	return "Ready"
}

func (r *OpsMateConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opsmatev1alpha1.OpsMateConfig{}).
		Complete(r)
}

func (r *OpsMateConfigReconciler) reconcileObject(ctx context.Context, desired client.Object) error {
	key := client.ObjectKeyFromObject(desired)

	switch desiredObject := desired.(type) {
	case *appsv1.Deployment:
		current := &appsv1.Deployment{}
		current.Name = key.Name
		current.Namespace = key.Namespace
		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, current, func() error {
			current.ObjectMeta.Labels = desiredObject.ObjectMeta.Labels
			current.ObjectMeta.OwnerReferences = desiredObject.ObjectMeta.OwnerReferences
			current.Spec = desiredObject.Spec
			return nil
		})
		return err
	case *corev1.Service:
		current := &corev1.Service{}
		current.Name = key.Name
		current.Namespace = key.Namespace
		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, current, func() error {
			current.ObjectMeta.Labels = desiredObject.ObjectMeta.Labels
			current.ObjectMeta.OwnerReferences = desiredObject.ObjectMeta.OwnerReferences
			current.Spec.Selector = desiredObject.Spec.Selector
			current.Spec.Ports = desiredObject.Spec.Ports
			return nil
		})
		return err
	case *corev1.ServiceAccount:
		current := &corev1.ServiceAccount{}
		current.Name = key.Name
		current.Namespace = key.Namespace
		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, current, func() error {
			current.ObjectMeta.Labels = desiredObject.ObjectMeta.Labels
			current.ObjectMeta.Annotations = desiredObject.ObjectMeta.Annotations
			current.ObjectMeta.OwnerReferences = desiredObject.ObjectMeta.OwnerReferences
			return nil
		})
		return err
	case *batchv1.Job:
		current := &batchv1.Job{}
		current.Name = key.Name
		current.Namespace = key.Namespace
		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, current, func() error {
			current.ObjectMeta.Labels = desiredObject.ObjectMeta.Labels
			current.ObjectMeta.OwnerReferences = desiredObject.ObjectMeta.OwnerReferences
			current.Spec = desiredObject.Spec
			return nil
		})
		return err
	default:
		current := desired.DeepCopyObject().(client.Object)
		current.SetName(key.Name)
		current.SetNamespace(key.Namespace)
		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, current, func() error {
			current.SetLabels(desired.GetLabels())
			current.SetOwnerReferences(desired.GetOwnerReferences())
			unstructuredCurrent, ok := current.(interface {
				UnstructuredContent() map[string]any
				SetUnstructuredContent(map[string]any)
			})
			if ok {
				unstructuredDesired := desired.(interface {
					UnstructuredContent() map[string]any
				})
				unstructuredCurrent.SetUnstructuredContent(unstructuredDesired.UnstructuredContent())
			}
			return nil
		})
		return err
	}
}

func upsertCondition(conditions []metav1.Condition, next metav1.Condition) []metav1.Condition {
	for index := range conditions {
		if conditions[index].Type == next.Type {
			conditions[index] = next
			return conditions
		}
	}
	return append(conditions, next)
}
