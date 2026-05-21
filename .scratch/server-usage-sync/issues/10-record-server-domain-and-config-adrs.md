Status: done

# 记录服务端命名与配置例外 ADR

## Type

HITL

## What to build

把两个首版服务端的关键架构取舍整理为 ADR，以便后续维护者理解为什么服务端核心域不沿用旧的代理语义，以及为什么首版配置层有意偏离 `server/CLAUDE.md` 中对 `viper` 的偏好。该切片需要产出能够被未来实现和评审直接引用的决策文档，而不是零散散落在聊天记录里的口头约定。完成后，后续重构或扩展时可以基于 ADR 判断这些决策是否仍然成立。

## Acceptance criteria

- [x] 记录“服务端核心域使用 `usage_reports` 而非 `proxy_request_logs`”的决策背景、取舍和影响
- [x] 记录“首版配置层使用标准库加环境变量而非 `viper`”的决策背景、取舍和影响
- [x] ADR 文档能够被未来实现者直接引用，并与当前 `CONTEXT.md` 保持一致
- [x] 该切片在进入执行前经过一次人工确认，确保 ADR 表述反映真实决策而不是实现猜测

## Blocked by

- `01-bootstrap-server-runtime-and-health.md`
- `02-persist-usage-reports-and-model-pricing.md`
