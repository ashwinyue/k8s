package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myappv1 "github.com/example/myapp-controller/pkg/apis/example/v1"
)

// MyAppReconciler 协调 MyApp 资源
type MyAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=myapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=myapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile 是核心的协调逻辑
func (r *MyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 获取 MyApp 实例
	myApp := &myappv1.MyApp{}
	err := r.Get(ctx, req.NamespacedName, myApp)
	if err != nil {
		if errors.IsNotFound(err) {
			// MyApp 资源已被删除
			logger.Info("MyApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// 获取失败
		logger.Error(err, "Failed to get MyApp")
		return ctrl.Result{}, err
	}

	// 更新状态为 Pending
	if myApp.Status.Phase == "" {
		logger.Info("Updating MyApp status to Pending", "name", myApp.Name, "namespace", myApp.Namespace)
		// 重新获取最新的 MyApp 对象以确保 ResourceVersion 正确
		latestMyApp := &myappv1.MyApp{}
		if err := r.Get(ctx, req.NamespacedName, latestMyApp); err != nil {
			logger.Error(err, "Failed to get latest MyApp for status update", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}

		logger.Info("Got latest MyApp for status update", "resourceVersion", latestMyApp.ResourceVersion)
		latestMyApp.Status.Phase = "Pending"
		latestMyApp.Status.Message = "正在创建资源"
		if err := r.Client.Status().Update(ctx, latestMyApp); err != nil {
			logger.Error(err, "Failed to update MyApp status", "name", latestMyApp.Name, "namespace", latestMyApp.Namespace, "resourceVersion", latestMyApp.ResourceVersion)
			// 如果状态更新失败，重新排队处理
			return ctrl.Result{RequeueAfter: time.Second * 5}, nil
		}
		logger.Info("Successfully updated MyApp status to Pending")
	}

	// 创建或更新 Deployment
	deployment := r.deploymentForMyApp(myApp)
	if err := ctrl.SetControllerReference(myApp, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	found := &appsv1.Deployment{}
	err = r.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
		// Deployment 创建成功，重新排队
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// 更新 Deployment 如果需要
	if found.Spec.Replicas != &myApp.Spec.Replicas {
		found.Spec.Replicas = &myApp.Spec.Replicas
		err = r.Update(ctx, found)
		if err != nil {
			logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// 规格更新成功，重新排队
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// 创建或更新 Service
	service := r.serviceForMyApp(myApp)
	if err := ctrl.SetControllerReference(myApp, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundService := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Create(ctx, service)
		if err != nil {
			logger.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		logger.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// 更新状态
	// 重新获取最新的 MyApp 对象以确保 ResourceVersion 正确
	finalMyApp := &myappv1.MyApp{}
	if err := r.Get(ctx, req.NamespacedName, finalMyApp); err != nil {
		logger.Error(err, "Failed to get latest MyApp for final status update")
		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	finalMyApp.Status.ReadyReplicas = found.Status.ReadyReplicas
	if found.Status.ReadyReplicas == myApp.Spec.Replicas {
		finalMyApp.Status.Phase = "Running"
		finalMyApp.Status.Message = "所有副本都已就绪"
	} else {
		finalMyApp.Status.Phase = "Pending"
		finalMyApp.Status.Message = fmt.Sprintf("等待副本就绪: %d/%d", found.Status.ReadyReplicas, myApp.Spec.Replicas)
	}

	if err := r.Client.Status().Update(ctx, finalMyApp); err != nil {
		logger.Error(err, "Failed to update MyApp status")
		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// deploymentForMyApp 为 MyApp 创建 Deployment
func (r *MyAppReconciler) deploymentForMyApp(m *myappv1.MyApp) *appsv1.Deployment {
	labels := map[string]string{
		"app": m.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &m.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: m.Spec.Image,
						Name:  "app",
						Ports: []corev1.ContainerPort{{
							ContainerPort: m.Spec.Port,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}
}

// serviceForMyApp 为 MyApp 创建 Service
func (r *MyAppReconciler) serviceForMyApp(m *myappv1.MyApp) *corev1.Service {
	labels := map[string]string{
		"app": m.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-service",
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt(int(m.Spec.Port)),
				Protocol:   corev1.ProtocolTCP,
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

// SetupWithManager 设置 Controller 与 Manager
func (r *MyAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.MyApp{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
