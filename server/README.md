# CC Usage Server

`cc-status/server` 是接收 Claude Code 使用记录、计算费用并提供查询统计 API 的服务端。

## 特性

- 完整的 REST API（认证、同步、统计、日志、模型定价）
- 使用 Go `embed` 打包前端静态资源，单一文件部署
- SQLite 数据持久化
- SPA 路由 fallback 支持

## 相关 ADR

- `docs/adr/0001-server-domain-uses-usage-reports.md`：说明为什么服务端核心域使用 `usage_reports`，而不是历史材料中的 `proxy_request_logs`
- `docs/adr/0002-server-config-uses-stdlib-and-env.md`：说明为什么首版配置层使用标准库与环境变量，而不是立即引入 `viper`

## 快速开始

### 方式1：使用构建好的单一二进制文件（推荐）

构建包含前端的完整版本：

```bash
# Windows
cd ..
.\build.ps1

# Linux/macOS
cd ..
make build

# 运行
$env:CC_USAGE_SERVER_AUTH_TOKEN="dev-token"   # Windows
export CC_USAGE_SERVER_AUTH_TOKEN="dev-token"  # Linux/macOS

./server/bin/server    # 或 server.exe on Windows
```

访问 http://localhost:8080 查看仪表板。

### 方式2：开发模式（不包含前端）

在仓库根目录执行：

```bash
export CC_USAGE_SERVER_AUTH_TOKEN="dev-token"
export CC_USAGE_SERVER_LISTEN_ADDR=":8080"
export CC_USAGE_SERVER_SQLITE_PATH="./server/data/server.db"
go run ./server/cmd/server
```

说明：

- `CC_USAGE_SERVER_AUTH_TOKEN` 必填
- `CC_USAGE_SERVER_LISTEN_ADDR` 默认 `:8080`
- `CC_USAGE_SERVER_SQLITE_PATH` 默认 `./server/data/server.db`

## 最小 Smoke 流程

健康检查：

```bash
curl http://127.0.0.1:8080/healthz
```

写入一条最小同步记录：

```bash
curl -X POST http://127.0.0.1:8080/api/v1/sync \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json" \
  -d @- <<'EOF'
{
  "client_id": "local-smoke",
  "reports": [
    {
      "request_id": "session:smoke-1",
      "app_type": "claude",
      "model": "claude-sonnet-4-0",
      "input_tokens": 100,
      "output_tokens": 200,
      "cache_read_tokens": 0,
      "cache_creation_tokens": 0,
      "created_at": 1743840000,
      "session_id": "smoke-session",
      "data_source": "session_log"
    }
  ]
}
EOF
```

查看原始日志：

```bash
curl -H "Authorization: Bearer dev-token" \
  "http://127.0.0.1:8080/api/v1/logs?client_id=local-smoke&limit=20"
```

查看总览和趋势：

```bash
curl -H "Authorization: Bearer dev-token" \
  http://127.0.0.1:8080/api/v1/stats/overview

curl -H "Authorization: Bearer dev-token" \
  "http://127.0.0.1:8080/api/v1/stats/trend?interval=day"
```

查看定价列表：

```bash
curl -H "Authorization: Bearer dev-token" \
  http://127.0.0.1:8080/api/v1/model-pricings
```

## API 接口

### 公开接口

- `GET /healthz`：健康检查，无需认证

### 需要认证的接口（Bearer Token）

- `POST /api/v1/sync`：接收客户端上报的使用记录
- `GET /api/v1/model-pricings`：获取模型定价列表
- `POST /api/v1/model-pricings`：创建模型定价
- `PUT /api/v1/model-pricings/:id`：更新模型定价
- `GET /api/v1/stats/overview`：获取统计总览
- `GET /api/v1/stats/trend`：获取趋势统计
- `GET /api/v1/stats/dashboard`：获取仪表盘数据
- `GET /api/v1/logs`：获取使用记录日志
- `GET /api/v1/ping`：ping 测试

## 与现有 Client 联调

在仓库根目录复制示例配置：

```bash
mkdir -p ~/.cc-usage-client
cp ./client/config.example.yaml ~/.cc-usage-client/config.yaml
```

把 `~/.cc-usage-client/config.yaml` 中的服务地址改成本地 server：

```yaml
server_url: "http://127.0.0.1:8080"
auth_token: "dev-token"
batch_size: 500
timeout_seconds: 30
```

然后执行：

```bash
go run ./client/cmd/cc-usage-client dry-run
go run ./client/cmd/cc-usage-client sync
```

联调后可以重新执行：

```bash
curl -H "Authorization: Bearer dev-token" \
  "http://127.0.0.1:8080/api/v1/logs?limit=20"

curl -H "Authorization: Bearer dev-token" \
  http://127.0.0.1:8080/api/v1/stats/overview
```

## 构建说明

Server 使用 Go `embed` 将 web 静态资源打包进二进制文件。

### 构建模式

- **开发模式**（默认）：不嵌入静态资源，仅提供 API
- **发布模式**（`-tags embed`）：嵌入 web 静态资源，同时提供 API 和 Web UI

### 构建命令

项目使用静态编译（`CGO_ENABLED=0`），生成的二进制文件不依赖系统 GLIBC，可以在 CentOS 7 等老系统上直接运行。

```bash
# 发布构建（推荐从项目根目录执行）
cd ..

# Windows
.\build.ps1

# Linux/macOS
make build
```

或在 server 目录单独构建：

```bash
cd server

# Linux/macOS: 静态编译
CGO_ENABLED=0 go build -tags embed -ldflags="-s -w -extldflags=-static" -o bin/server ./cmd/server

# Windows: 静态编译
$env:CGO_ENABLED=0
go build -tags embed -ldflags="-s -w" -o bin/server.exe ./cmd/server
```

### 目录结构

```
server/
├── cmd/server/         # 入口
├── internal/
│   ├── handler/       # HTTP handler + UI 文件服务
│   ├── service/       # 业务逻辑
│   ├── repository/    # 数据访问
│   ├── model/         # 数据模型
│   ├── middleware/    # Gin 中间件
│   └── config/        # 配置
├── data/              # SQLite 数据目录
└── internal/ui/       # UI 构建产物（不提交到 git）
```

## 开发验证

在仓库根目录执行：

```bash
go test ./server/...
```
