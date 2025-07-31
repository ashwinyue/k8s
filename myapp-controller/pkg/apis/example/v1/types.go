package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MyAppSpec 定义 MyApp 的期望状态
type MyAppSpec struct {
	// Image 容器镜像
	Image string `json:"image"`
	// Replicas 副本数量
	Replicas int32 `json:"replicas"`
	// Port 服务端口
	Port int32 `json:"port"`
}

// MyAppStatus 定义 MyApp 的实际状态
type MyAppStatus struct {
	// Phase 应用阶段
	Phase string `json:"phase,omitempty"`
	// Message 状态消息
	Message string `json:"message,omitempty"`
	// ReadyReplicas 就绪的副本数
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MyApp 是我们自定义资源的定义
type MyApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyAppSpec   `json:"spec,omitempty"`
	Status MyAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// MyAppList 包含 MyApp 的列表
type MyAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyApp `json:"items"`
}
