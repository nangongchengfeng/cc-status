Status: done

# 交付 Claude dry-run 扫描链路

## What to build

交付一个可验证的 `dry-run` 端到端链路：扫描 Claude Code 的**会话日志**，发现主会话 JSONL 和 `subagents` 文件，按既有 Rust 规则提取**Token 使用记录**，并输出统计结果。该切片必须复用真实解析逻辑，包括 assistant 过滤、`stop_reason` 过滤、`output_tokens > 0` 过滤、单文件内 `message.id` 去重，以及尾部半写入坏行的瞬时错误处理。`dry-run` 允许初始化基础依赖，但不能推进**同步状态**或写入**已上报记录**。

## Acceptance criteria

- [x] `dry-run` 能扫描 `~/.claude/projects/` 下的主会话与 `subagents` JSONL 文件，并输出扫描文件数与提取到的**Token 使用记录**统计
- [x] 提取规则严格满足 assistant、`message.stop_reason`、`output_tokens > 0` 与单文件内重复 `message.id` 合并要求
- [x] 遇到尾部半写入或 JSON 解析失败时，只阻断当前文件并报告错误，不会推进该文件的**同步状态**
- [x] `dry-run` 可初始化数据库和**客户端标识**，但不会写入或更新 `sync_state` 与 `reported_ids`

## Blocked by

- [01-bootstrap-client-runtime-and-state.md](./01-bootstrap-client-runtime-and-state.md)
