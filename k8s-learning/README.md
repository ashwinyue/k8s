# Kubernetes 学习教程

欢迎来到 Kubernetes 学习之旅！这个目录包含了从基础到实践的完整学习材料。

## 📚 教程目录

### 1. [基础概念](./01-k8s-basics.md)
- Kubernetes 简介
- 核心概念（Pod、Deployment、Service 等）
- 基本命令
- 学习路径规划

### 2. [简单 Pod 示例](./02-simple-pod.yaml)
- 第一个 Pod YAML 配置
- 资源限制设置
- 基本操作命令

### 3. [Deployment 和 Service](./03-deployment-example.yaml)
- 多副本管理
- 服务发现和负载均衡
- 滚动更新和扩缩容

### 4. [配置管理](./04-configmap-secret.yaml)
- ConfigMap 使用方法
- Secret 安全存储
- 环境变量和文件挂载

### 5. [实践指南](./05-practice-guide.md)
- 环境搭建步骤
- 逐步实践教程
- 故障排查技巧
- 进阶学习方向

## 🚀 快速开始

1. **环境准备**
   ```bash
   # 安装 kubectl 和 minikube
   brew install kubectl minikube
   
   # 启动本地集群
   minikube start
   ```

2. **验证环境**
   ```bash
   kubectl cluster-info
   kubectl get nodes
   ```

3. **开始学习**
   - 阅读 [基础概念](./01-k8s-basics.md)
   - 按照 [实践指南](./05-practice-guide.md) 逐步操作

## 📋 学习检查清单

- [ ] 理解 Kubernetes 基本概念
- [ ] 成功创建和管理 Pod
- [ ] 掌握 Deployment 的使用
- [ ] 配置 Service 进行服务发现
- [ ] 使用 ConfigMap 和 Secret 管理配置
- [ ] 学会基本的故障排查
- [ ] 完成所有实践练习

## 🛠️ 常用命令速查

```bash
# 查看资源
kubectl get pods
kubectl get deployments
kubectl get services

# 应用配置
kubectl apply -f <yaml-file>

# 查看详情
kubectl describe <resource-type> <resource-name>

# 查看日志
kubectl logs <pod-name>

# 进入容器
kubectl exec -it <pod-name> -- /bin/bash
```

## 🔗 相关资源

- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [kubectl 命令参考](https://kubernetes.io/docs/reference/kubectl/)
- [YAML 配置参考](https://kubernetes.io/docs/reference/)
- [Minikube 文档](https://minikube.sigs.k8s.io/docs/)

## 💡 学习建议

1. **循序渐进**：按照文件编号顺序学习
2. **动手实践**：每个概念都要亲自操作
3. **理解原理**：不仅要知道怎么做，还要知道为什么
4. **多做实验**：尝试修改配置，观察结果
5. **记录笔记**：记录遇到的问题和解决方案

## 🤝 获得帮助

如果在学习过程中遇到问题：
1. 查看 [实践指南](./05-practice-guide.md) 中的故障排查部分
2. 使用 `kubectl describe` 和 `kubectl logs` 命令调试
3. 查阅官方文档
4. 在社区论坛寻求帮助

祝你学习顺利！🎉