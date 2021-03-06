package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/go-logr/logr"
	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	"github.com/rflorenc/drwatcher-operator/controllers"
	drwversion "github.com/rflorenc/drwatcher-operator/version"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	veleroInstall "github.com/vmware-tanzu/velero/pkg/install"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(drv1.AddToScheme(scheme))
	utilruntime.Must(velerov1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

var (
	metricsAddr string
	metricsPort string = ":8080"
)

func checkPreRequisites(ctx context.Context, client client.Client, logger logr.Logger) error {
	for _, unstructuredCrd := range veleroInstall.AllCRDs().Items {
		crd := &apiextv1beta1.CustomResourceDefinition{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredCrd.Object, crd)
		if err != nil {
			logger.Error(err, "Required velero.io group CRDs could not be converted.")
			return err
		}
		err = client.Get(ctx, types.NamespacedName{Name: crd.ObjectMeta.Name}, crd)
		if err != nil {
			if errors.IsNotFound(err) {
				logger.Error(err, fmt.Sprintf("Required velero.io CRDs %v not found", crd.Name))
				return err
			}
		}
	}
	return nil
}

func main() {
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", metricsPort, "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	ctx := context.Background()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "5683d144.drx",
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	if err = (&controllers.DRWatcherReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DRWatcher"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DRWatcher")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	err = checkPreRequisites(ctx, mgr.GetClient(), setupLog)
	if err != nil {
		setupLog.Error(err, "unable to find existing Velero installation", "controller", "DRWatcher")
		os.Exit(1)
	}

	setupLog.Info(fmt.Sprintf("Operator Version: %s", drwversion.Version))
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
