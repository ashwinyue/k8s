package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion 是组版本，用于注册这些对象
var SchemeGroupVersion = schema.GroupVersion{Group: "example.com", Version: "v1"}

// Kind 获取给定对象的 GroupVersionKind
func Kind(kind string) schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(kind)
}

// Resource 获取给定 GroupVersion 的 GroupVersionResource
func Resource(resource string) schema.GroupVersionResource {
	return SchemeGroupVersion.WithResource(resource)
}

var (
	// SchemeBuilder 初始化一个 SchemeBuilder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme 将这个组版本的类型添加到给定的 scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// addKnownTypes 将已知类型添加到 scheme
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&MyApp{},
		&MyAppList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
