package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DRWatcherSpec defines the desired state of DRWatcher
type DRWatcherSpec struct {
	// Important: Run "make manifests" to regenerate code after modifying this file

	Schedule       string `json:"schedule,omitempty"`
	Command        string `json:"command,omitempty"`
	ReadyForBackup bool   `json:"readyForBackup,omitempty"`
}

// DRWatcherStatus defines the observed state of DRWatcher
type DRWatcherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase string `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DRWatcher is the Schema for the drwatchers API
type DRWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DRWatcherSpec   `json:"spec,omitempty"`
	Status DRWatcherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DRWatcherList contains a list of DRWatcher
type DRWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DRWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DRWatcher{}, &DRWatcherList{})
}
