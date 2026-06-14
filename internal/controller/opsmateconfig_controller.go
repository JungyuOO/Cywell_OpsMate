package controller

import (
	"context"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/appserver"
	consoleplugin "github.com/JungyuOO/Cywell_OpsMate/internal/controller/console"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller/postgres"
	appsv1 "k8s.io/api/apps/v1"
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
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=console.openshift.io,resources=consoleplugins,verbs=get;list;watch;create;update;patch
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

	config.Status.OverallStatus = "Ready"
	ready := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		ObservedGeneration: config.Generation,
		LastTransitionTime: metav1.Now(),
		Reason:             "ResourcesReconciled",
		Message:            "Appserver and PostgreSQL resources are reconciled.",
	}
	config.Status.Conditions = upsertCondition(config.Status.Conditions, ready)
	if err := r.Status().Update(ctx, config); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
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
