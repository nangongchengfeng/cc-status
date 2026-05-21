Status: done

# 完善 CLI 反馈与核心回归测试

## What to build

为首版 client 补齐面向操作者的结果反馈与最小但关键的自动化回归测试。该切片聚焦两个目标：一是让 `dry-run` 和 `sync` 的输出足够清楚，能直接说明扫描文件数、提取记录数、上报成功数、跳过数和错误摘要；二是把解析、状态存储和同步编排的核心外部行为固化成跨平台自动化测试，降低后续修改回归风险。

## Acceptance criteria

- [x] `dry-run` 和 `sync` 都会输出清晰、稳定的结果摘要，覆盖扫描文件数、提取记录数、上报成功数、跳过数和错误数量
- [x] 自动化测试至少覆盖 Claude 日志解析、SQLite 状态存储、同步编排三个核心模块，并聚焦外部行为而非内部实现细节
- [x] 测试包含尾部坏行、单文件内重复 `message.id`、成功推进 `sync_state`、失败批次不推进状态、`dry-run` 不写业务状态等关键回归场景
- [x] 相关测试可通过 Go 标准测试框架在本地独立运行，不依赖真实中心服务

## Blocked by

- [02-dry-run-claude-log-scan.md](./02-dry-run-claude-log-scan.md)
- [04-failure-recovery-and-partial-commit.md](./04-failure-recovery-and-partial-commit.md)
