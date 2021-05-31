package controllers

import (
	"context"

	"github.com/go-logr/logr"
	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *DRWatcherReconciler) newRestoreForCR(cr *drv1.DRWatcher) *velerov1.Restore {
	var includedNamespaces []string
	includedNamespaces = append(includedNamespaces, cr.Namespace)

	labels := map[string]string{
		"created-by":                 "drwatcher",
		"velero.io/storage-location": "default",
	}

	return &velerov1.Restore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.BackupName,
			Namespace: veleroNamespace,
			Labels:    labels,
		},
		Spec: velerov1.RestoreSpec{
			BackupName:         cr.Spec.BackupName,
			IncludedNamespaces: includedNamespaces,
		},
	}
}

func (r *DRWatcherReconciler) getRestoreNames(ctx context.Context, cr *drv1.DRWatcher, logger logr.Logger) []string {
	restoreList := &velerov1.RestoreList{}
	restoreListOptions := []client.ListOption{
		client.InNamespace(veleroNamespace),
	}
	err := r.List(ctx, restoreList, restoreListOptions...)
	if err != nil {
		logger.Error(err, "Failed to list restores", "Namespace",
			veleroNamespace, "DRWatcher.Name", cr.Name)
	}
	var restoreInfo []string
	for _, backup := range restoreList.Items {
		restoreInfo = append(restoreInfo, backup.Name)
	}
	return restoreInfo
}
