package controllers

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var testClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))
	// KUBEBUILDER_ASSETS can also be used
	// https://book.kubebuilder.io/reference/envtest.html#configuring-your-test-control-plane
	Expect(os.Setenv("TEST_ASSET_KUBE_APISERVER", "/tmp/testbin/kube-apiserver")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_ETCD", "/tmp/testbin/etcd")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_KUBECTL", "/tmp/testbin/kubectl")).To(Succeed())

	By("bootstrapping test environment")
	t := true
	if os.Getenv("TEST_USE_EXISTING_CLUSTER") == "true" {
		testEnv = &envtest.Environment{
			UseExistingCluster: &t,
		}
	} else {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:        []string{filepath.Join("..", "config", "crd", "bases")},
			AttachControlPlaneOutput: true,
			ControlPlane:             testEnv.ControlPlane,
			ControlPlaneStartTimeout: 360,
			ControlPlaneStopTimeout:  360,
			UseExistingCluster:       &t,
			// KubeAPIServerFlags: append(
			// 	envtest.DefaultKubeAPIServerFlags,
			// 	"--advertise-address=127.0.0.1",
			// ),
		}
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = drv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	testManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).NotTo(HaveOccurred())

	err = (&DRWatcherReconciler{
		Client: testManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("drwatcher"),
	}).SetupWithManager(testManager)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		err = testManager.Start(ctrl.SetupSignalHandler())
		Expect(err).NotTo(HaveOccurred())
	}()

	testClient = testManager.GetClient()
	Expect(testClient).ToNot(BeNil())

	close(done)
}, 10000)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(os.Unsetenv("TEST_ASSET_KUBE_APISERVER")).To(Succeed())
	Expect(os.Unsetenv("TEST_ASSET_ETCD")).To(Succeed())
	Expect(os.Unsetenv("TEST_ASSET_KUBECTL")).To(Succeed())
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
