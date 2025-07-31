# Kubernetes 实践指南

## 环境准备

### 1. 安装必要工具

```bash
# 安装 kubectl (macOS)
brew install kubectl

# 安装 minikube (本地 Kubernetes 集群)
brew install minikube

# 或者使用 Docker Desktop 的 Kubernetes 功能
```

### 2. 启动本地集群

```bash
# 启动 minikube
minikube start

# 验证集群状态
kubectl cluster-info
kubectl get nodes
```

## 实践步骤

### 第一步：基础 Pod 操作

1. 应用简单的 Pod：
   ```bash
   kubectl apply -f 02-simple-pod.yaml
   ```

2. 观察 Pod 状态：
   ```bash
   kubectl get pods
   kubectl get pods -o wide  # 显示更多信息
   ```

3. 查看 Pod 详情：
   ```bash
   kubectl describe pod nginx-pod
   ```

4. 查看 Pod 日志：
   ```bash
   kubectl logs nginx-pod
   ```

5. 进入 Pod 容器：
   ```bash
   kubectl exec -it nginx-pod -- /bin/bash
   # 在容器内执行：curl localhost
   # 退出：exit
   ```

6. 删除 Pod：
   ```bash
   kubectl delete -f 02-simple-pod.yaml
   ```

### 第二步：Deployment 和 Service

1. 应用 Deployment 和 Service：
   ```bash
   kubectl apply -f 03-deployment-example.yaml
   ```

2. 观察资源创建：
   ```bash
   kubectl get deployments
   kubectl get pods
   kubectl get services
   ```

3. 测试扩缩容：
   ```bash
   # 扩展到 5 个副本
   kubectl scale deployment nginx-deployment --replicas=5
   kubectl get pods
   
   # 缩减到 2 个副本
   kubectl scale deployment nginx-deployment --replicas=2
   kubectl get pods
   ```

4. 测试滚动更新：
   ```bash
   # 更新镜像版本
   kubectl set image deployment/nginx-deployment nginx=nginx:1.21
   
   # 观察滚动更新过程
   kubectl rollout status deployment/nginx-deployment
   kubectl get pods -w  # 实时观察
   ```

5. 测试服务访问：
   ```bash
   # 端口转发到本地
   kubectl port-forward service/nginx-service 8080:80
   # 在另一个终端测试：curl http://localhost:8080
   ```

### 第三步：ConfigMap 和 Secret

1. 应用配置资源：
   ```bash
   kubectl apply -f 04-configmap-secret.yaml
   ```

2. 查看配置：
   ```bash
   kubectl get configmaps
   kubectl get secrets
   kubectl describe configmap app-config
   ```

3. 验证环境变量：
   ```bash
   kubectl exec -it app-pod -- env | grep -E "DATABASE_URL|LOG_LEVEL|USERNAME"
   ```

4. 查看挂载的文件：
   ```bash
   kubectl exec -it app-pod -- ls /etc/config
   kubectl exec -it app-pod -- cat /etc/config/app.properties
   ```

### 第四步：故障排查

1. 查看事件：
   ```bash
   kubectl get events --sort-by=.metadata.creationTimestamp
   ```

2. 查看资源使用：
   ```bash
   kubectl top nodes  # 需要 metrics-server
   kubectl top pods
   ```

3. 调试 Pod：
   ```bash
   kubectl describe pod <pod-name>
   kubectl logs <pod-name> --previous  # 查看之前容器的日志
   ```

## 清理资源

```bash
# 删除所有创建的资源
kubectl delete -f 02-simple-pod.yaml
kubectl delete -f 03-deployment-example.yaml
kubectl delete -f 04-configmap-secret.yaml

# 或者删除整个命名空间（如果使用了自定义命名空间）
# kubectl delete namespace <namespace-name>
```

## 下一步学习

1. **持久化存储**：学习 PersistentVolume 和 PersistentVolumeClaim
2. **网络**：学习 Ingress 控制器和网络策略
3. **安全**：学习 RBAC、ServiceAccount 和 Pod 安全策略
4. **监控**：学习 Prometheus 和 Grafana 集成
5. **CI/CD**：学习与 Jenkins、GitLab CI 等工具集成

## 有用的命令速查

```bash
# 快速创建资源
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port=80 --type=NodePort

# 生成 YAML 模板
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml

# 查看 API 资源
kubectl api-resources

# 查看 API 版本
kubectl api-versions

# 设置默认命名空间
kubectl config set-context --current --namespace=<namespace>
```

祝你学习愉快！🚀