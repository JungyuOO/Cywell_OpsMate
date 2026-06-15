package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// OpsMateConfigSpec defines the desired state for the Cywell OpsMate Operator.
type OpsMateConfigSpec struct {
	Lightspeed LightspeedSpec `json:"lightspeed"`
	AIOps      AIOpsSpec      `json:"aiops,omitempty"`
	RAG        []RAGSpec      `json:"rag,omitempty"`
	Embedding  EmbeddingSpec  `json:"embedding,omitempty"`
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

type EmbeddingSpec struct {
	EndpointURL          string `json:"endpointURL,omitempty"`
	Model                string `json:"model,omitempty"`
	Dimensions           int    `json:"dimensions,omitempty"`
	CredentialsSecretRef string `json:"credentialsSecretRef,omitempty"`
	CredentialsSecretKey string `json:"credentialsSecretKey,omitempty"`
	RequirePGVector      bool   `json:"requirePGVector,omitempty"`
	RetrievalMode        string `json:"retrievalMode,omitempty"`
	RetrievalSlowMillis  int    `json:"retrievalSlowMillis,omitempty"`
}

type DatabaseSpec struct {
	Type                      string `json:"type,omitempty"`
	Image                     string `json:"image,omitempty"`
	DSNSecretRef              string `json:"dsnSecretRef,omitempty"`
	DSNSecretKey              string `json:"dsnSecretKey,omitempty"`
	PGVectorMigrationApproved bool   `json:"pgVectorMigrationApproved,omitempty"`
	SharedBuffers             string `json:"sharedBuffers,omitempty"`
	MaxConnections            int    `json:"maxConnections,omitempty"`
}

type ConsoleSpec struct {
	Enabled                       bool     `json:"enabled,omitempty"`
	DisplayName                   string   `json:"displayName,omitempty"`
	AdminTokenSecretRef           string   `json:"adminTokenSecretRef,omitempty"`
	AdminTokenSecretKey           string   `json:"adminTokenSecretKey,omitempty"`
	AdminUsers                    []string `json:"adminUsers,omitempty"`
	AdminGroups                   []string `json:"adminGroups,omitempty"`
	AdminAuthProxyEnabled         bool     `json:"adminAuthProxyEnabled,omitempty"`
	AdminAuthProxyImage           string   `json:"adminAuthProxyImage,omitempty"`
	AdminAuthProxyCookieSecretRef string   `json:"adminAuthProxyCookieSecretRef,omitempty"`
	AdminRouteHost                string   `json:"adminRouteHost,omitempty"`
}

type OpsMateConfigStatus struct {
	OverallStatus     string             `json:"overallStatus,omitempty"`
	PGVectorReady     bool               `json:"pgVectorReady,omitempty"`
	PGVectorLastError string             `json:"pgVectorLastError,omitempty"`
	Reembedding       ReembeddingStatus  `json:"reembedding,omitempty"`
	Conditions        []metav1.Condition `json:"conditions,omitempty"`
}

type ReembeddingStatus struct {
	Running   bool   `json:"running,omitempty"`
	Processed int    `json:"processed,omitempty"`
	Failed    int    `json:"failed,omitempty"`
	LastError string `json:"lastError,omitempty"`
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
