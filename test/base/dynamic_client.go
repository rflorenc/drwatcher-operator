package base

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const timeout = time.Second * 30

// DynamicWaitForDeletion uses the dynamic client to wait for a resource to not exist.
func DynamicWaitForDeletion(dynamicClient dynamic.Interface, gvr schema.GroupVersionResource, namespace, name string, logger log.FieldLogger) error {
	logr := logger.WithFields(log.Fields{
		"gvr":       gvr,
		"namespace": namespace,
		"name":      name,
	})

	for start := time.Now(); ; {
		if time.Since(start) > timeout {
			logr.Error("resource not deleted before timeout")
			return fmt.Errorf("resource not deleted before timeout")
		}
		_, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			logr.Info("resource successfully deleted")
			return nil
		} else if err != nil {
			logr.WithError(err).Info("unexpected error getting resource")
		} else {
			logr.Info("resource still exists, sleeping...")
		}

		time.Sleep(500 * time.Millisecond)
	}
}
