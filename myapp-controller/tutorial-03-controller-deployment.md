# 第三课：Controller 部署和测试

## 项目结构

```
myapp-controller/
├── cmd/
│   └── main.go                    # 主程序入口
├── pkg/
│   ├── apis/example/v1/
│   │   ├── types.go               # 资源类型定义
│   │   ├── register.go            # Scheme 注册
│   │   └── zz_generated.deepcopy.go # DeepCopy 方法
│   └── controller/
│       └── myapp_controller.go    # Controller 逻辑
├── config/
│   ├── crd/
│   │   └── 01-simple-crd.yaml     # CRD 定义
│   ├── rbac/
│   │   └── rbac.yaml              # RBAC 配置
│   └── manager/
│       └── deployment.yaml        # Controller 部署
├── Dockerfile                     # 容器镜像构建
├── Makefile                       # 构建脚本
├── go.mod                         # Go 模块
└── go.sum                         # 依赖校验
```

## 本地开发和测试

### 1. 构建项目
```bash
# 格式化代码
make fmt

# 检查代码
make vet

# 构建二进制文件
make build
```

### 2. 本地运行 Controller
```bash
# 确保 CRD 已安装
kubectl apply -f config/crd/01-simple-crd.yaml

# 本地运行 Controller
make run
```

### 3. 测试 Controller
在另一个终端中：
```bash
# 创建测试资源
kubectl apply -f ../02-myapp-example.yaml

# 观察 Controller 日志
# 查看创建的 Deployment 和 Service
kubectl get deployments
kubectl get services
kubectl get myapps

# 查看 MyApp 状态
kubectl describe myapp nginx-app
```

## 容器化部署

### 1. 构建镜像
```bash
# 构建 Docker 镜像
make docker-build

# 加载镜像到 kind 集群
kind load docker-image myapp-controller:latest
```

### 2. 部署到集群
```bash
# 创建 RBAC
kubectl apply -f config/rbac/rbac.yaml

# 部署 Controller
kubectl apply -f config/manager/deployment.yaml

# 检查部署状态
kubectl get pods -n myapp-system
kubectl logs -n myapp-system deployment/myapp-controller
```

## Controller 工作原理

### 1. 监听机制
Controller 使用 Informer 监听以下资源的变化：
- MyApp 资源（主要监听对象）
- Deployment 资源（拥有的资源）
- Service 资源（拥有的资源）

### 2. 协调循环
当资源发生变化时，Controller 执行以下步骤：

1. **获取 MyApp 资源**
   - 如果资源不存在，说明已被删除，结束处理
   - 如果获取失败，记录错误并重试

2. **更新状态为 Pending**
   - 首次创建时设置状态

3. **创建或更新 Deployment**
   - 根据 MyApp 规格创建 Deployment
   - 设置 OwnerReference 确保级联删除
   - 如果 Deployment 不存在则创建
   - 如果副本数不匹配则更新

4. **创建或更新 Service**
   - 为应用创建 ClusterIP Service
   - 设置 OwnerReference

5. **更新 MyApp 状态**
   - 根据 Deployment 状态更新 MyApp 状态
   - 设置就绪副本数和阶段信息

### 3. 错误处理和重试
- 使用 `ctrl.Result{RequeueAfter: time.Minute}` 定期重新协调
- 错误时返回错误，触发指数退避重试
- 使用结构化日志记录操作和错误

## 高级特性

### 1. Owner Reference
```go
if err := ctrl.SetControllerReference(myApp, deployment, r.Scheme); err != nil {
    return ctrl.Result{}, err
}
```
- 确保当 MyApp 被删除时，相关的 Deployment 和 Service 也会被删除
- 建立资源间的父子关系

### 2. 状态管理
```go
myApp.Status.ReadyReplicas = found.Status.ReadyReplicas
if found.Status.ReadyReplicas == myApp.Spec.Replicas {
    myApp.Status.Phase = "Running"
} else {
    myApp.Status.Phase = "Pending"
}
```
- 实时反映应用的实际状态
- 提供用户友好的状态信息

### 3. RBAC 权限
Controller 需要以下权限：
- 对 MyApp 资源的完全访问权限
- 对 Deployment 和 Service 的 CRUD 权限
- 对 MyApp 状态子资源的更新权限

## 故障排查

### 1. Controller 无法启动
```bash
# 检查 RBAC 权限
kubectl auth can-i create deployments --as=system:serviceaccount:myapp-system:myapp-controller

# 检查 CRD 是否安装
kubectl get crd myapps.example.com

# 查看 Controller 日志
kubectl logs -n myapp-system deployment/myapp-controller
```

### 2. 资源未创建
```bash
# 检查 MyApp 资源状态
kubectl describe myapp <name>

# 查看 Controller 日志中的错误
kubectl logs -n myapp-system deployment/myapp-controller | grep ERROR
```

### 3. 状态不更新
```bash
# 检查 Controller 是否有状态更新权限
kubectl auth can-i update myapps/status --as=system:serviceaccount:myapp-system:myapp-controller
```

## 练习任务

1. **修改 Controller 逻辑**
   - 添加对容器环境变量的支持
   - 实现资源限制的配置

2. **增强状态报告**
   - 添加更详细的错误信息
   - 实现事件记录

3. **添加验证**
   - 实现 Admission Webhook
   - 添加资源验证逻辑

## 下一步

在下一课中，我们将学习：
- Operator 模式和最佳实践
- 使用 Operator SDK 快速开发
- 实现复杂的生命周期管理
- 添加监控和告警