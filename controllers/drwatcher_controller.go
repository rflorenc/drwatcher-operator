package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	backupVolumesAnnotation string = "backup.velero.io/backup-volumes"
	veleroNamespace         string = "velero"
)

// DRWatcherReconciler reconciles a DRWatcher object
type DRWatcherReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile tracks changes to DRWatcher CRs and enables self service creation of Backups and Schedules.
// +kubebuilder:rbac:groups=dr.seven,resources=drwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dr.seven,resources=drwatchers/status,verbs=get;update;patch
func (r *DRWatcherReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("DRWatcher", req.NamespacedName)
	var err error
	var drwatcherCR drv1.DRWatcher

	// get the DRWatcher CR
	err = r.Get(ctx, req.NamespacedName, &drwatcherCR)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("[-] watcher instance not found.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	resticAnnotations, err := r.getPodInfo(ctx, &drwatcherCR, logger)
	if err != nil {
		logger.Error(err, "Failed to get restic pod annotations", "Namespace",
			veleroNamespace, "DRWatcher.Name", drwatcherCR.Name)
		return ctrl.Result{}, err
	} else {
		logger.Info(fmt.Sprintf("resticAnnotations: %s", resticAnnotations))
	}

	if drwatcherCR.Spec.ReadyForBackup {
		if drwatcherCR.Spec.Schedule != "" {
			existingSchedules := r.getScheduleNames(ctx, &drwatcherCR, logger)
			logger.Info(fmt.Sprintf("Current schedules: %s", existingSchedules))

			newSchedule := r.newScheduleForCR(&drwatcherCR)
			logger.Info(fmt.Sprintf("Creating new Schedule %s with Spec.Schedule %s for project: %s", newSchedule.Name, newSchedule.Spec.Schedule, drwatcherCR.Namespace))

			err = r.Create(ctx, newSchedule)
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to create Schedule %s with Spec.Schedule %s for project: %s", newSchedule.Name, newSchedule.Spec.Schedule, drwatcherCR.Namespace))
				return ctrl.Result{}, err
			}
		} else {
			existingBackups := r.getBackupNames(ctx, &drwatcherCR, logger)
			logger.Info(fmt.Sprintf("Current backups: %s", existingBackups))

			newBackup := r.newBackupForCR(&drwatcherCR)
			logger.Info(fmt.Sprintf("Creating backup %s for project: %s", newBackup.Name, drwatcherCR.Namespace))

			err = r.Create(ctx, newBackup)
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to create backup %s for project: %s", newBackup.Name, drwatcherCR.Namespace))
				return ctrl.Result{}, err
			}
		}
	} else {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up a new DRWatcher controller managed by mgr
func (r *DRWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&drv1.DRWatcher{}).
		Owns(&corev1.Pod{}).
		// Owns(&velerov1.Backup{}).
		Complete(r)
}
