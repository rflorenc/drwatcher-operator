package controllers

import (
	"context"

	"github.com/go-logr/logr"
	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *DRWatcherReconciler) newBackupForCR(cr *drv1.DRWatcher) *velerov1.Backup {
	var includedNamespaces []string
	includedNamespaces = append(includedNamespaces, cr.Namespace)

	labels := map[string]string{
		"created-by":                 "drwatcher",
		"velero.io/storage-location": "default",
	}

	return &velerov1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Namespace + "-drwatcher-backup",
			Namespace: veleroNamespace,
			Labels:    labels,
		},
		Spec: velerov1.BackupSpec{
			StorageLocation:    "default",
			IncludedNamespaces: includedNamespaces,
		},
	}
}

func (r *DRWatcherReconciler) getBackupNames(ctx context.Context, cr *drv1.DRWatcher, logger logr.Logger) []string {
	backupList := &velerov1.BackupList{}
	backupListOptions := []client.ListOption{
		client.InNamespace(veleroNamespace),
	}
	if err := r.List(ctx, backupList, backupListOptions...); err != nil {
		logger.Error(err, "Failed to list backups", "Namespace",
			veleroNamespace, "DRWatcher.Name", cr.Name)
	}
	var backupInfo []string
	for _, backup := range backupList.Items {
		backupInfo = append(backupInfo, backup.Name)
	}
	return backupInfo
}
