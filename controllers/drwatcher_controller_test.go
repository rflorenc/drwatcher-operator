package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
)

var _ = Describe("DRWatcher controller", func() {
	typeMeta := metav1.TypeMeta{
		APIVersion: "dr.test/v1",
		Kind:       "DRWatcher",
	}

	objectMeta := metav1.ObjectMeta{
		Name:      "test-drwatcher",
		Namespace: "default",
	}

	spec := drv1.DRWatcherSpec{
		Schedule:       "",
		Command:        "test",
		ReadyForBackup: false,
	}

	const (
		timeout  = time.Minute
		interval = time.Millisecond * 250
	)

	Context("When creating DRWatcher instance", func() {
		It("should return an error if Name is not present", func(done Done) {
			ctx := context.Background()
			drwatcher := &drv1.DRWatcher{
				TypeMeta:   typeMeta,
				ObjectMeta: objectMeta,
				Spec:       spec,
			}
			Expect(testClient.Create(ctx, drwatcher)).Should(Succeed())

			key := types.NamespacedName{Name: objectMeta.Name, Namespace: objectMeta.Namespace}
			obj := &drv1.DRWatcher{}

			Eventually(func() bool {
				err := testClient.Get(ctx, key, obj)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(obj.ObjectMeta.Name).ShouldNot(Equal(""))
		})
	})
})
