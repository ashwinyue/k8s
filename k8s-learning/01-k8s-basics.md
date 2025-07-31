# Kubernetes 基础教程

## 什么是 Kubernetes？

Kubernetes（简称 K8s）是一个开源的容器编排平台，用于自动化部署、扩展和管理容器化应用程序。

## 核心概念

### 1. Pod
- Kubernetes 中最小的部署单元
- 包含一个或多个容器
- 共享网络和存储

### 2. Deployment
- 管理 Pod 的副本数量
- 提供滚动更新和回滚功能
- 确保应用的高可用性

### 3. Service
- 为 Pod 提供稳定的网络访问
- 负载均衡流量到多个 Pod
- 支持不同类型：ClusterIP、NodePort、LoadBalancer

### 4. Namespace
- 提供资源隔离
- 组织和管理集群资源
- 支持多租户环境

### 5. ConfigMap 和 Secret
- ConfigMap：存储非敏感配置数据
- Secret：存储敏感信息（密码、密钥等）

## 基本命令

```bash
# 查看集群信息
kubectl cluster-info

# 查看节点
kubectl get nodes

# 查看 Pod
kubectl get pods

# 查看所有资源
kubectl get all

# 创建资源
kubectl apply -f <yaml-file>

# 删除资源
kubectl delete -f <yaml-file>

# 查看资源详情
kubectl describe <resource-type> <resource-name>

# 查看日志
kubectl logs <pod-name>

# 进入容器
kubectl exec -it <pod-name> -- /bin/bash
```

## 下一步

1. 安装 kubectl 和 minikube（本地开发）
2. 创建第一个 Pod
3. 学习 YAML 配置文件
4. 实践 Deployment 和 Service
5. 探索更高级的概念（Ingress、PersistentVolume 等）

## 实践练习

接下来我们将创建一些实际的 YAML 文件来练习这些概念。