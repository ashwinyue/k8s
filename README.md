# Kubernetes 二次开发学习教程

## 环境信息
- Kind集群: kind
- Kubernetes版本: v1.33.1
- kubectl版本: v1.32.3

## 学习路径

### 第一阶段：CRD (Custom Resource Definition) 基础
1. 理解CRD概念
2. 创建简单的CRD
3. 使用kubectl操作自定义资源

### 第二阶段：Controller 开发
1. 理解Controller模式
2. 使用client-go开发Controller
3. 实现资源的CRUD操作

### 第三阶段：Operator 开发
1. 理解Operator模式
2. 使用Operator SDK
3. 实现复杂的业务逻辑

### 第四阶段：高级特性
1. Webhook开发
2. 自定义调度器
3. 扩展API Server

## 开始第一个实例：创建CRD

让我们从最基础的CRD开始学习Kubernetes二次开发。