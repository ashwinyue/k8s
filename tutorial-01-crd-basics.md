# 第一课：CRD 基础教程

## 什么是 CRD？

Custom Resource Definition (CRD) 是 Kubernetes 提供的扩展机制，允许用户定义自己的资源类型。通过 CRD，我们可以像使用内置资源（如 Pod、Service）一样使用自定义资源。

## CRD 的核心概念

### 1. Group（组）
- 用于组织相关的 API 资源
- 例如：`example.com`

### 2. Version（版本）
- API 的版本控制
- 例如：`v1`, `v1beta1`

### 3. Kind（类型）
- 资源的类型名称
- 例如：`MyApp`

### 4. Resource（资源）
- 资源的复数形式名称
- 例如：`myapps`

## 实践步骤

### 步骤 1：应用 CRD
```bash
# 创建 CRD
kubectl apply -f 01-simple-crd.yaml

# 验证 CRD 是否创建成功
kubectl get crd myapps.example.com

# 查看 CRD 详细信息
kubectl describe crd myapps.example.com
```

### 步骤 2：创建自定义资源
```bash
# 创建自定义资源实例
kubectl apply -f 02-myapp-example.yaml

# 查看创建的资源
kubectl get myapps
kubectl get ma  # 使用短名称

# 查看详细信息
kubectl describe myapp nginx-app
```

### 步骤 3：操作自定义资源
```bash
# 编辑资源
kubectl edit myapp nginx-app

# 删除资源
kubectl delete myapp redis-app

# 使用 YAML 输出
kubectl get myapp nginx-app -o yaml
```

## CRD Schema 详解

### OpenAPI v3 Schema
我们的 CRD 使用 OpenAPI v3 Schema 来定义资源的结构：

```yaml
schema:
  openAPIV3Schema:
    type: object
    properties:
      spec:          # 期望状态
        type: object
        properties:
          image:       # 容器镜像
            type: string
          replicas:    # 副本数
            type: integer
            minimum: 1
            maximum: 10
          port:        # 端口
            type: integer
            minimum: 1
            maximum: 65535
      status:          # 实际状态
        type: object
        properties:
          phase:       # 阶段
            type: string
            enum: ["Pending", "Running", "Failed"]
          message:     # 消息
            type: string
```

### Additional Printer Columns
自定义 `kubectl get` 命令的输出列：

```yaml
additionalPrinterColumns:
- name: Image
  type: string
  jsonPath: .spec.image
- name: Replicas
  type: integer
  jsonPath: .spec.replicas
- name: Status
  type: string
  jsonPath: .status.phase
- name: Age
  type: date
  jsonPath: .metadata.creationTimestamp
```

## 练习任务

1. 修改 CRD，添加一个新的字段 `env`（环境变量）
2. 创建一个新的 MyApp 资源实例
3. 尝试创建一个不符合 schema 的资源，观察错误信息
4. 使用 `kubectl patch` 命令更新资源

## 下一步

现在你已经了解了 CRD 的基础知识，下一课我们将学习如何开发 Controller 来监听和处理这些自定义资源的变化。