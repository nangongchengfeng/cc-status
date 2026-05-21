# CC Usage Server

`cc-status/server` 是接收 Claude Code 使用记录、计算费用并提供查询统计 API 的服务端。

## 本地启动

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

## 最小 Smoke 流程

健康检查：

```powershell
curl http://127.0.0.1:8080/healthz
```

写入一条最小同步记录：

```powershell
$headers = @{
  Authorization = "Bearer dev-token"
  "Content-Type" = "application/json"
}

$body = @'
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
'@

Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:8080/api/v1/sync" -Headers $headers -Body $body
```

查看原始日志：

```powershell
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/logs?client_id=local-smoke&limit=20" -Headers $headers
```

查看总览和趋势：

```powershell
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/stats/overview" -Headers $headers
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/stats/trend?interval=day" -Headers $headers
```

查看定价列表：

```powershell
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/model-pricings" -Headers $headers
```

## 与现有 Client 联调

在 `client` 目录复制示例配置：

```powershell
New-Item -ItemType Directory -Force -Path "$HOME/.cc-usage-client" | Out-Null
Copy-Item .\client\config.example.yaml "$HOME/.cc-usage-client/config.yaml"
```

把 `~/.cc-usage-client/config.yaml` 中的服务地址改成本地 server：

```yaml
server_url: "http://127.0.0.1:8080"
auth_token: "dev-token"
batch_size: 500
timeout_seconds: 30
```

然后执行：

```powershell
go run ./client/cmd/cc-usage-client dry-run
go run ./client/cmd/cc-usage-client sync
```

联调后可以重新执行：

```powershell
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/logs?limit=20" -Headers $headers
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/api/v1/stats/overview" -Headers $headers
```

## 开发验证

在仓库根目录执行：

```powershell
go test ./server/...
```
