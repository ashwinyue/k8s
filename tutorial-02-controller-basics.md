# 第二课：Controller 基础教程

## 什么是 Controller？

Controller 是 Kubernetes 的核心概念，它负责监听资源的变化并执行相应的操作，确保集群的实际状态与期望状态一致。

## Controller 模式

Controller 遵循以下模式：
1. **Watch（监听）**：监听特定资源的变化
2. **Compare（比较）**：比较当前状态与期望状态
3. **Act（执行）**：执行必要的操作来达到期望状态

## 开发环境准备

### 1. 安装 Go
```bash
# 检查 Go 版本（需要 1.19+）
go version
```

### 2. 初始化 Go 模块
```bash
# 创建项目目录
mkdir myapp-controller
cd myapp-controller

# 初始化 Go 模块
go mod init github.com/example/myapp-controller
```

### 3. 安装依赖
```bash
# 安装 client-go 和相关依赖
go get k8s.io/client-go@latest
go get k8s.io/apimachinery@latest
go get k8s.io/api@latest
go get sigs.k8s.io/controller-runtime@latest
```

## Controller 架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Server    │◄──►│   Controller    │◄──►│   Work Queue    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       │                       │
         │                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Informer    │    │   Reconciler    │    │   Event Handler │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 核心组件

### 1. Client
- 与 Kubernetes API Server 通信
- 执行 CRUD 操作

### 2. Informer
- 监听资源变化
- 本地缓存
- 事件通知

### 3. Work Queue
- 事件队列
- 重试机制
- 限流控制

### 4. Reconciler
- 核心业务逻辑
- 状态协调
- 错误处理

## 实现步骤

### 步骤 1：生成客户端代码
我们需要为自定义资源生成客户端代码。有两种方式：

#### 方式 1：使用 code-generator
```bash
# 安装 code-generator
go install k8s.io/code-generator/cmd/client-gen@latest
go install k8s.io/code-generator/cmd/lister-gen@latest
go install k8s.io/code-generator/cmd/informer-gen@latest
```

#### 方式 2：使用 controller-runtime（推荐）
controller-runtime 提供了更简单的方式来开发 Controller。

### 步骤 2：定义资源类型
创建 Go 结构体来表示我们的自定义资源。

### 步骤 3：实现 Reconciler
实现核心的协调逻辑。

### 步骤 4：启动 Controller
配置并启动 Controller。

## 下一步实践

在下一个教程中，我们将：
1. 创建完整的 Controller 项目结构
2. 实现 MyApp Controller
3. 处理 MyApp 资源的创建、更新和删除
4. 创建对应的 Deployment 和 Service

## 学习资源

- [client-go 官方文档](https://github.com/kubernetes/client-go)
- [controller-runtime 文档](https://github.com/kubernetes-sigs/controller-runtime)
- [Kubernetes Controller 模式](https://kubernetes.io/docs/concepts/architecture/controller/)

## 练习思考

1. Controller 和 Operator 有什么区别？
2. 为什么需要 Work Queue？
3. Informer 的本地缓存有什么作用？
4. 如何处理 Controller 的错误和重试？