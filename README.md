# CC Status

`cc-status` 是一个围绕 Claude Code 使用量采集与汇总的 Go 项目，当前包含三个独立子模块：

- `client`：单次执行 CLI，从本地 Claude Code JSONL 会话日志中提取 token 使用记录并增量上报
- `server`：集中接收上报、计算费用、持久化使用记录，并提供查询与管理 API
- `web`：React 前端仪表板，展示使用量统计与趋势

项目当前以本地 SQLite 为默认状态存储，适合先完成单机闭环与联调验证。





## 仓库结构

```text
.
├── client/              # Go CLI：扫描 Claude Code 日志并上报
├── server/              # Go API：接收、计费、查询、管理
├── web/                 # React 前端：仪表板展示
├── docs/adr/           # 架构决策记录
├── .scratch/           # 本地 PRD / issues 追踪
├── CONTEXT.md          # 当前领域语言与约束
└── claude_usage_stats.md # 历史背景与字段语义参考
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
- `GET /api/v1/stats/dashboard`
- `GET /api/v1/logs`
- 除 `GET /healthz` 外，全部 API 需要 Bearer token

更多说明见：[server/README.md](./server/README.md)

### Web

- 费用概览卡片
- 费用趋势与 Token 趋势图表
- 模型与客户端排行
- 缓存分析
- 最近请求列表

更多说明见：[web/README.md](./web/README.md)

## 快速开始

### 1. 启动服务端

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

### 2. 准备客户端配置

在仓库根目录执行：

```bash
mkdir -p ~/.cc-usage-client
cp ./client/config.example.yaml ~/.cc-usage-client/config.yaml
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

```bash
go run ./client/cmd/cc-usage-client dry-run
go run ./client/cmd/cc-usage-client sync
```

### 4. 查询服务端结果

```bash
# 健康检查
curl http://127.0.0.1:8080/healthz

# 查看日志
curl -H "Authorization: Bearer dev-token" \
  "http://127.0.0.1:8080/api/v1/logs?limit=20"

# 查看总览
curl -H "Authorization: Bearer dev-token" \
  http://127.0.0.1:8080/api/v1/stats/overview

# 查看趋势
curl -H "Authorization: Bearer dev-token" \
  "http://127.0.0.1:8080/api/v1/stats/trend?interval=day"
```

### 5. 启动前端

在另一个终端执行：

```bash
cd web
npm install
npm run dev
```

访问 `http://localhost:5173` 查看仪表板。

## 开发验证

本仓库包含三个独立模块，请分别进入子目录执行测试：

```bash
# Client 测试
cd client
go test ./...

# Server 测试
cd ../server
go test ./...

# Web 测试
cd ../web
npm run test
```

## 文档入口

- 领域语言与约束：[CONTEXT.md](./CONTEXT.md)
- Client 使用说明：[client/README.md](./client/README.md)
- Server 使用说明：[server/README.md](./server/README.md)
- Web 使用说明：[web/README.md](./web/README.md)
- ADR：[docs/adr/](./docs/adr/)

## 本地 PRD / Issue 追踪

项目当前使用本地 markdown 追踪需求与 issue：

- `client` 侧：`.scratch/claude-usage-sync/`
- `server` 侧：`.scratch/server-usage-sync/`

其中：

- `PRD.md` 用于记录需求与设计
- `issues/*.md` 用于记录拆分任务与状态
