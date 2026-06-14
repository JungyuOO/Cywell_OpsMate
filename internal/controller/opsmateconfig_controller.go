package controller

import (
	"context"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
// +kubebuilder:rbac:groups=opsmate.cywell.io,resources=opsmateconfigs/status,verbs=get
func (r *OpsMateConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *OpsMateConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opsmatev1alpha1.OpsMateConfig{}).
		Complete(r)
}
