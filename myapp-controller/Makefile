# 镜像 URL，用于构建或推送镜像目标
IMG ?= myapp-controller:latest

# 获取当前运行的操作系统架构信息
ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)

# 设置 Shell
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ 通用

# 显示帮助信息
.PHONY: help
help: ## 显示帮助信息
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ 开发

.PHONY: fmt
fmt: ## 运行 go fmt
	go fmt ./...

.PHONY: vet
vet: ## 运行 go vet
	go vet ./...

.PHONY: test
test: fmt vet ## 运行测试
	go test ./... -coverprofile cover.out

##@ 构建

.PHONY: build
build: fmt vet ## 构建 manager 二进制文件
	go build -o bin/manager cmd/main.go

.PHONY: run
run: fmt vet ## 在本地运行 controller
	go run ./cmd/main.go

.PHONY: docker-build
docker-build: test ## 构建 docker 镜像
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## 推送 docker 镜像
	docker push ${IMG}

##@ 部署

.PHONY: install
install: ## 安装 CRD 到集群
	kubectl apply -f config/crd/

.PHONY: uninstall
uninstall: ## 从集群卸载 CRD
	kubectl delete -f config/crd/

.PHONY: deploy
deploy: ## 部署 controller 到集群
	kubectl apply -f config/manager/

.PHONY: undeploy
undeploy: ## 从集群卸载 controller
	kubectl delete -f config/manager/

##@ 构建依赖

.PHONY: deps
deps: ## 下载依赖
	go mod download
	go mod tidy