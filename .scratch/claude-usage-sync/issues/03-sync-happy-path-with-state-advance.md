Status: done

# 交付 sync 成功上报链路

## What to build

在 `dry-run` 的扫描与解析能力之上，交付 `sync` 的成功路径：把新提取的**Token 使用记录**按固定批大小发送到预留的中心服务接口，并在服务端明确确认整批成功后写入**已上报记录**与对应文件的**同步状态**。该切片重点验证“成功即推进”的主流程，包括 `request_id` 生成、批大小默认值 500、HTTP 请求组装、超时时间设置和无新增数据时的正常退出。

## Acceptance criteria

- [x] `sync` 会为每条**Token 使用记录**生成 `request_id = session:<message_id>`，并按最多 500 条一批发送到配置中的服务端 URL
- [x] 当服务端成功确认整批记录时，client 会写入对应 `reported_ids`，并推进相关文件的 `sync_state`
- [x] 无新增记录时，`sync` 会输出清晰结果并以 0 退出
- [x] 成功上报路径会携带稳定的**客户端标识**与约定好的请求字段结构，便于后续 server 对接

## Blocked by

- [01-bootstrap-client-runtime-and-state.md](./01-bootstrap-client-runtime-and-state.md)
- [02-dry-run-claude-log-scan.md](./02-dry-run-claude-log-scan.md)
