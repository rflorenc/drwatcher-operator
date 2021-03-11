package controllers

import (
	"context"

	"github.com/go-logr/logr"
	drv1 "github.com/rflorenc/drwatcher-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *DRWatcherReconciler) getPodInfo(ctx context.Context, cr *drv1.DRWatcher, logger logr.Logger) ([]string, error) {
	var err error
	podList := &corev1.PodList{}
	podListOptions := []client.ListOption{
		client.InNamespace(cr.Namespace),
	}

	if err = r.List(ctx, podList, podListOptions...); err != nil {
		logger.Error(err, "Failed to list pods", "DRWatcher.Namespace",
			cr.Namespace, "DRWatcher.Name", cr.Name)
		return nil, err
	}

	podResticAnnotations := getPodAnnotations(podList.Items)
	if podResticAnnotations != nil {
		return podResticAnnotations, nil
	}
	return nil, err
}

func getPodAnnotations(pods []corev1.Pod) []string {
	var podInfo []string
	for _, pod := range pods {
		podInfo = append(podInfo, pod.GetAnnotations()[backupVolumesAnnotation])
	}
	return podInfo
}
