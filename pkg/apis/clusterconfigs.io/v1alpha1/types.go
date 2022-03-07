package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterConfigMap struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec ClusterConfigMapSpec `json:"spec,omitempty"`
}

type ClusterConfigMapSpec struct {
	// +optional
	Data map[string]string `json:"data,omitempty"`

	// +optional
	BinaryData map[string][]byte `json:"binaryData,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterConfigMapList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ClusterConfigMap `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterSecret struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec ClusterSecretSpec `json:"spec,omitempty"`
}

type ClusterSecretSpec struct {
	// +optional
	Data map[string][]byte `json:"data,omitempty"`

	// +optional
	StringData map[string]string `json:"stringData,omitempty"`

	// +optional
	Type corev1.SecretType `json:"type,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterSecretList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ClusterSecret `json:"items"`
}
