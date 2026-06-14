package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// OpsMateConfigSpec defines the desired state for the Cywell OpsMate Operator.
type OpsMateConfigSpec struct {
	Lightspeed LightspeedSpec `json:"lightspeed"`
	AIOps      AIOpsSpec      `json:"aiops,omitempty"`
	RAG        []RAGSpec      `json:"rag,omitempty"`
	Database   DatabaseSpec   `json:"database,omitempty"`
	Console    ConsoleSpec    `json:"console,omitempty"`
}

type LightspeedSpec struct {
	APIBaseURL           string `json:"apiBaseURL,omitempty"`
	CredentialsSecretRef string `json:"credentialsSecretRef,omitempty"`
	DefaultProvider      string `json:"defaultProvider,omitempty"`
	DefaultModel         string `json:"defaultModel,omitempty"`
}

type AIOpsSpec struct {
	Enabled bool `json:"enabled,omitempty"`
}

type RAGSpec struct {
	Image     string `json:"image,omitempty"`
	IndexPath string `json:"indexPath,omitempty"`
	IndexID   string `json:"indexID,omitempty"`
}

type DatabaseSpec struct {
	Type           string `json:"type,omitempty"`
	SharedBuffers  string `json:"sharedBuffers,omitempty"`
	MaxConnections int    `json:"maxConnections,omitempty"`
}

type ConsoleSpec struct {
	Enabled     bool   `json:"enabled,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type OpsMateConfigStatus struct {
	OverallStatus string             `json:"overallStatus,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type OpsMateConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpsMateConfigSpec   `json:"spec,omitempty"`
	Status OpsMateConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type OpsMateConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []OpsMateConfig `json:"items"`
}
