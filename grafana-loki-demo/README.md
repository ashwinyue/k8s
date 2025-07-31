# Grafana Loki 日志采集体系

这是一个完整的 Grafana Loki 日志采集和可视化系统，包含以下组件：

## 🏗️ 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Demo App   │───▶│    Promtail     │───▶│      Loki       │
│   (日志生成)     │    │   (日志采集)     │    │   (日志存储)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │    Grafana      │
                                               │   (日志可视化)   │
                                               └─────────────────┘
```

## 📦 组件说明

### 1. Loki
- **作用**: 日志聚合系统，类似 Prometheus 但专门用于日志
- **特点**: 不索引日志内容，只索引标签，存储成本低
- **端口**: 3100

### 2. Grafana
- **作用**: 数据可视化平台，提供日志查询和仪表板
- **特点**: 支持多种数据源，强大的查询和可视化能力
- **端口**: 3000
- **默认账号**: admin/admin

### 3. Promtail
- **作用**: 日志采集代理，负责收集和转发日志到 Loki
- **特点**: 轻量级，支持多种日志格式和来源
- **端口**: 9080

### 4. Go Demo App
- **作用**: 演示应用，生成各种类型的日志
- **特点**: 包含结构化日志、错误日志、访问日志等
- **端口**: 8080

## 🚀 快速开始

### 1. 启动服务
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 2. 访问服务
- **Grafana**: http://localhost:3000 (admin/admin)
- **Loki**: http://localhost:3100
- **Go Demo App**: http://localhost:8080
- **Promtail**: http://localhost:9080

### 3. 配置 Grafana
1. 登录 Grafana
2. 添加 Loki 数据源: http://loki:3100
3. 导入预配置的仪表板
4. 开始查询日志

## 📁 目录结构

```
grafana-loki-demo/
├── README.md                 # 项目说明
├── docker-compose.yml        # Docker 编排文件
├── configs/                  # 配置文件目录
│   ├── loki.yml             # Loki 配置
│   ├── promtail.yml         # Promtail 配置
│   └── grafana/             # Grafana 配置
│       ├── datasources/     # 数据源配置
│       └── dashboards/      # 仪表板配置
├── go-demo/                 # Go 演示应用
│   ├── main.go             # 主程序
│   ├── go.mod              # Go 模块文件
│   ├── Dockerfile          # Docker 镜像
│   └── logs/               # 日志输出目录
└── data/                   # 数据持久化目录
    ├── loki/               # Loki 数据
    └── grafana/            # Grafana 数据
```

## 🔍 日志查询示例

### LogQL 查询语法
```logql
# 查看所有日志
{job="go-demo"}

# 查看错误日志
{job="go-demo"} |= "ERROR"

# 查看特定时间范围的日志
{job="go-demo"}[5m]

# 统计错误日志数量
sum(count_over_time({job="go-demo"} |= "ERROR"[5m]))

# 按级别分组统计
sum by (level) (count_over_time({job="go-demo"}[5m]))
```

## 🛠️ 开发指南

### 添加新的日志源
1. 修改 `promtail.yml` 添加新的 scrape_configs
2. 重启 Promtail 服务
3. 在 Grafana 中创建新的查询

### 自定义仪表板
1. 在 Grafana 中创建新仪表板
2. 导出 JSON 配置
3. 保存到 `configs/grafana/dashboards/`

## 📊 监控指标

- 日志摄入速率
- 错误日志比例
- 服务响应时间
- 系统资源使用情况

## 🔧 故障排查

### 常见问题
1. **Loki 无法启动**: 检查配置文件语法和端口占用
2. **Promtail 无法采集**: 检查日志文件路径和权限
3. **Grafana 无法连接 Loki**: 检查网络连接和数据源配置

### 日志查看
```bash
# 查看服务日志
docker-compose logs loki
docker-compose logs promtail
docker-compose logs grafana
docker-compose logs go-demo
```

## 📚 参考资料

- [Loki 官方文档](https://grafana.com/docs/loki/)
- [Grafana 官方文档](https://grafana.com/docs/grafana/)
- [Promtail 配置指南](https://grafana.com/docs/loki/latest/clients/promtail/)
- [LogQL 查询语言](https://grafana.com/docs/loki/latest/logql/)