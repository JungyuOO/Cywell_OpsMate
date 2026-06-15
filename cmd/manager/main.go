package main

import (
	"os"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	"github.com/JungyuOO/Cywell_OpsMate/internal/controller"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opsmatev1alpha1.AddToScheme(scheme))

	options := ctrl.Options{
		Scheme: scheme,
	}
	if namespace := watchedNamespace(); namespace != "" {
		options.Cache.DefaultNamespaces = map[string]cache.Config{
			namespace: {},
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		ctrl.Log.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := controller.SetupOpsMateConfigReconciler(mgr); err != nil {
		ctrl.Log.Error(err, "unable to create OpsMateConfig controller")
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "manager exited")
		os.Exit(1)
	}
}

func watchedNamespace() string {
	if namespace := os.Getenv("WATCH_NAMESPACE"); namespace != "" {
		return namespace
	}
	return os.Getenv("POD_NAMESPACE")
}
