# CC Status

`cc-status` 是一个围绕 Claude Code 使用量采集与汇总的 Go 项目，当前包含两个独立子模块：

- `client`：单次执行 CLI，从本地 Claude Code JSONL 会话日志中提取 token 使用记录并增量上报
- `server`：集中接收上报、计算费用、持久化使用记录，并提供查询与管理 API

项目当前以本地 SQLite 为默认状态存储，适合先完成单机闭环与联调验证。

## 仓库结构

```text
.
├─ client/                     # Go CLI：扫描 Claude Code 日志并上报
├─ server/                     # Go API：接收、计费、查询、管理
├─ docs/adr/                   # 架构决策记录
├─ .scratch/                   # 本地 PRD / issues 追踪
├─ CONTEXT.md                  # 当前领域语言与约束
└─ claude_usage_stats.md       # 历史背景与字段语义参考
```

## 当前能力

### Client

- `dry-run`：扫描本地日志并输出摘要，不推进业务同步状态
- `sync`：批量上报新记录，成功后写入 `reported_ids` 并推进 `sync_state`
- 默认配置目录：`~/.cc-usage-client`
- 默认本地状态库：`~/.cc-usage-client/client.db`

更多说明见：[client/README.md](./client/README.md)

### Server

- `POST /api/v1/sync`：接收上报批次，兼容旧 client 成功响应格式
- `GET /api/v1/model-pricings` / `POST /api/v1/model-pricings` / `PUT /api/v1/model-pricings/:id`
- `GET /api/v1/stats/overview`
- `GET /api/v1/stats/trend`
- `GET /api/v1/logs`
- 除 `GET /healthz` 外，全部 API 需要 Bearer token

更多说明见：[server/README.md](./server/README.md)

## 快速开始

### 1. 启动服务端

在仓库根目录执行：

```powershell
$env:CC_USAGE_SERVER_AUTH_TOKEN = "dev-token"
$env:CC_USAGE_SERVER_LISTEN_ADDR = ":8080"
$env:CC_USAGE_SERVER_SQLITE_PATH = ".\server\data\server.db"
go run ./server/cmd/server
```

说明：

- `CC_USAGE_SERVER_AUTH_TOKEN` 必填
- `CC_USAGE_SERVER_LISTEN_ADDR` 默认 `:8080`
- `CC_USAGE_SERVER_SQLITE_PATH` 默认 `./server/data/server.db`

### 2. 准备客户端配置

在仓库根目录执行：

```powershell
New-Item -ItemType Directory -Force -Path "$HOME/.cc-usage-client" | Out-Null
Copy-Item .\client\config.example.yaml "$HOME/.cc-usage-client/config.yaml"
```

将 `~/.cc-usage-client/config.yaml` 至少修改为：

```yaml
server_url: "http://127.0.0.1:8080"
auth_token: "dev-token"
batch_size: 500
timeout_seconds: 30
```

### 3. 预演与同步

在仓库根目录执行：

```powershell
go run ./client/cmd/cc-usage-client dry-run
go run ./client/cmd/cc-usage-client sync
```

### 4. 查询服务端结果

```powershell
$headers = @{ Authorization = "Bearer dev-token" }

Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/healthz"
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/logs?limit=20" -Headers $headers
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/stats/overview" -Headers $headers
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/stats/trend?interval=day" -Headers $headers
```

## 开发验证

本仓库当前是两个独立 Go module，请分别进入子目录执行测试：

```powershell
Set-Location .\client
go test ./...
```

```powershell
Set-Location .\server
go test ./...
```

## 文档入口

- 领域语言与约束：[CONTEXT.md](./CONTEXT.md)
- Client 使用说明：[client/README.md](./client/README.md)
- Server 使用说明：[server/README.md](./server/README.md)
- ADR：[0001-server-domain-uses-usage-reports.md](./docs/adr/0001-server-domain-uses-usage-reports.md)
- ADR：[0002-server-config-uses-stdlib-and-env.md](./docs/adr/0002-server-config-uses-stdlib-and-env.md)

## 本地 PRD / Issue 追踪

项目当前使用本地 markdown 追踪需求与 issue：

- `client` 侧：`.scratch/claude-usage-sync/`
- `server` 侧：`.scratch/server-usage-sync/`

其中：

- `PRD.md` 用于记录需求与设计
- `issues/*.md` 用于记录拆分任务与状态
