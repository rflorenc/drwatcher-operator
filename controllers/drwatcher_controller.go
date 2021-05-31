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

// Reconcile tracks changes to DRWatcher CRs and enables self service creation of Backups, Restores and Schedules.
// +kubebuilder:rbac:groups=dr.seven,resources=drwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dr.seven,resources=drwatchers/status,verbs=get;update;patch
func (r *DRWatcherReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("DRWatcher", req.NamespacedName)
	var err error
	var drwatcherCR drv1.DRWatcher

	if drwatcherCR.Spec.ReadyForRestore && drwatcherCR.Spec.ReadyForBackup {
		return ctrl.Result{}, err
	}

	err = r.Get(ctx, req.NamespacedName, &drwatcherCR)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "watcher instance not found.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	resticAnnotations, err := r.getPodInfo(ctx, &drwatcherCR, logger)
	if err != nil {
		logger.Error(err, "Failed to get restic pod annotations", "Namespace",
			veleroNamespace, "DRWatcher.Name", drwatcherCR.Name)
		return ctrl.Result{}, err
	}
	logger.Info(fmt.Sprintf("resticAnnotations: %s", resticAnnotations))

	if drwatcherCR.Spec.ReadyForBackup {
		if drwatcherCR.Spec.Schedule != "" {
			existingSchedules := r.getScheduleNames(ctx, &drwatcherCR, logger)
			logger.Info(fmt.Sprintf("Current schedules: %s", existingSchedules))

			newSchedule := r.newScheduleForCR(&drwatcherCR)
			logger.Info(fmt.Sprintf("Creating new Schedule %s with Spec.Schedule %s for project: %s", newSchedule.Name, newSchedule.Spec.Schedule, drwatcherCR.Namespace))

			err = r.Create(ctx, newSchedule)
			if err != nil {
				logger.Error(err, fmt.Sprintf("Failed to create Schedule %s with Spec.Schedule %s for project: %s", newSchedule.Name, newSchedule.Spec.Schedule, drwatcherCR.Namespace))
				return ctrl.Result{}, err
			}
		} else {
			existingBackups := r.getBackupNames(ctx, &drwatcherCR, logger)
			logger.Info(fmt.Sprintf("Current backups: %s", existingBackups))

			newBackup := r.newBackupForCR(&drwatcherCR)
			logger.Info(fmt.Sprintf("Creating backup %s for project: %s", newBackup.Name, drwatcherCR.Namespace))

			err = r.Create(ctx, newBackup)
			if err != nil {
				logger.Error(err, fmt.Sprintf("Failed to create backup %s for project: %s", newBackup.Name, drwatcherCR.Namespace))
				return ctrl.Result{}, err
			}
		}
	}

	if drwatcherCR.Spec.ReadyForRestore {
		drwatcherCR.Spec.ReadyForBackup = false

		existingRestores := r.getRestoreNames(ctx, &drwatcherCR, logger)
		logger.Info(fmt.Sprintf("Current restores: %s", existingRestores))

		newRestore := r.newRestoreForCR(&drwatcherCR)
		logger.Info(fmt.Sprintf("Creating restore %s for project: %s", newRestore.Name, drwatcherCR.Namespace))

		err = r.Create(ctx, newRestore)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Failed to create restore %s for project: %s", newRestore.Name, drwatcherCR.Namespace))
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up a new DRWatcher controller managed by mgr
func (r *DRWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&drv1.DRWatcher{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
