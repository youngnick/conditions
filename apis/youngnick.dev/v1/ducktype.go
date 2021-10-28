package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DuckConditions is a duck-type for a Runtime.Object with a status.Conditions field.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type DuckConditions struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Status DuckStatus `json:"status,omitempty"`
}

// DuckStatus holds the Conditions struct.
type DuckStatus struct {
	// +optional
	//
	// If you are another controller owner and wish to add a condition, you *should*
	// namespace your condition with a label, like `controller.domain.com/ConditionName`.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}
