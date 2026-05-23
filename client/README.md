# Claude Usage Sync Client

基于 Go 的单次执行 CLI，用于从本地 Claude Code 会话日志中提取 token 使用数据，并同步到中心服务。

当前已实现：

- `dry-run`：扫描本地日志并输出摘要，不写入业务同步状态
- `sync`：扫描、批量上报、记录 `reported_ids`，并在满足条件时推进 `sync_state`

## 目录与数据位置

默认应用数据目录：

```text
~/.cc-usage-client
```

该目录下会保存：

- `config.yaml`：运行配置
- `client.db`：本地 SQLite 数据库
- `client.lock`：进程级互斥锁文件

默认 Claude Code 日志目录：

```text
~/.claude/projects
```

支持扫描：

- 主会话 `*.jsonl`
- 子 Agent 会话 `subagents/*.jsonl`

## 构建与运行

在仓库根目录执行：

```bash
# 构建
cd client
go build ./cmd/cc-usage-client

# 或者直接运行
go run ./cmd/cc-usage-client dry-run
go run ./cmd/cc-usage-client sync
```

本地预演扫描：

```bash
go run ./cmd/cc-usage-client dry-run
```

正式同步：

```bash
go run ./cmd/cc-usage-client sync
```

## 配置文件

程序默认从 `~/.cc-usage-client/config.yaml` 读取配置。

可以先复制仓库里的示例模板：

```bash
mkdir -p ~/.cc-usage-client
cp ./client/config.example.yaml ~/.cc-usage-client/config.yaml
```

配置字段说明：

- `server_url`：中心服务地址，不带 `/api/v1/sync` 也可以，程序会自动拼接固定路径
- `auth_token`：Bearer Token，建议优先用环境变量覆盖
- `batch_size`：单批上报记录数，默认 `500`
- `timeout_seconds`：单次 HTTP 请求超时秒数，默认 `30`
- `client_name`：客户端名称（可选），设置后会替代自动生成的 UUID，用于在服务端更友好地标识客户端

## 环境变量覆盖

以下环境变量会覆盖 `config.yaml` 中的同名配置：

- `CC_USAGE_CLIENT_SERVER_URL`
- `CC_USAGE_CLIENT_AUTH_TOKEN`
- `CC_USAGE_CLIENT_BATCH_SIZE`
- `CC_USAGE_CLIENT_TIMEOUT_SECONDS`
- `CC_USAGE_CLIENT_CLIENT_NAME`

示例：

```bash
export CC_USAGE_CLIENT_SERVER_URL="https://usage.example.com"
export CC_USAGE_CLIENT_AUTH_TOKEN="your-token"
go run ./client/cmd/cc-usage-client sync
```

## 输出样例

`dry-run` 成功：

```text
dry-run summary: files_scanned=12 records=37 errors=0
```

`dry-run` 遇到坏尾行：

```text
dry-run summary: files_scanned=12 records=36 errors=1
dry-run error: /path/to/session.jsonl: parse json line: unexpected end of JSON input
```

`sync` 成功：

```text
sync summary: files_scanned=12 records=37 accepted=37 skipped=0 errors=0
```

`sync` 无新增记录：

```text
sync summary: files_scanned=12 records=0 accepted=0 skipped=0 errors=0
```

其中：

- `files_scanned`：本次扫描到的 `.jsonl` 文件数
- `records`：本次实际准备上报的记录数
- `accepted`：服务端确认接收的记录数
- `skipped`：已上报或服务端返回重复的记录数
- `errors`：本次汇总的文件级错误数

## 当前同步规则

Claude 日志解析规则与现有设计保持一致：

- 仅处理 `type == "assistant"` 的消息
- 必须存在 `message.stop_reason`
- 必须满足 `output_tokens > 0`
- 同一文件内若 `message.id` 重复：
  - 优先保留有 `stop_reason` 的记录
  - 若都满足，则保留 `output_tokens` 更大的记录

`request_id` 生成规则：

```text
session:<message_id>
```

## 失败恢复语义

- 网络错误和 HTTP `5xx` 最多重试 3 次
- HTTP `200` 但业务码异常，整批按失败处理
- `accepted_count + duplicate_count` 小于本批记录数，整批按失败处理
- 成功批次会保留 `reported_ids`
- 失败批次不会写入 `reported_ids`
- 某个文件只有在本轮所有未上报记录都成功后，才推进该文件的 `sync_state`
- 文件被截断或替换后，会按从头重扫处理

## 开发验证

在 `client` 目录执行：

```bash
cd client
go test ./...
```
