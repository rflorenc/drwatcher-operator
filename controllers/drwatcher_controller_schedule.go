package controllers

import (
	"context"

	"github.com/go-logr/logr"
	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *DRWatcherReconciler) newScheduleForCR(cr *drv1.DRWatcher) *velerov1.Schedule {
	var includedNamespaces []string
	includedNamespaces = append(includedNamespaces, cr.Namespace)

	labels := map[string]string{
		"created-by": "drwatcher",
	}

	return &velerov1.Schedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Namespace + "-drwatcher-schedule",
			Namespace: veleroNamespace,
			Labels:    labels,
		},
		Spec: velerov1.ScheduleSpec{
			Schedule: cr.Spec.Schedule,
			Template: velerov1.BackupSpec{
				IncludedNamespaces: includedNamespaces,
			},
		},
	}
}

func (r *DRWatcherReconciler) getScheduleNames(ctx context.Context, cr *drv1.DRWatcher, logger logr.Logger) []string {
	// list existing velero schedule CRs in the velero namespace
	scheduleList := &velerov1.ScheduleList{}
	scheduleListOptions := []client.ListOption{
		client.InNamespace(veleroNamespace),
	}
	if err := r.List(ctx, scheduleList, scheduleListOptions...); err != nil {
		logger.Error(err, "Failed to list schedules", "Namespace",
			veleroNamespace, "DRWatcher.Name", cr.Name)
	}

	var ScheduleInfo []string
	for _, schedule := range scheduleList.Items {
		ScheduleInfo = append(ScheduleInfo, schedule.Name)
	}
	return ScheduleInfo
}
