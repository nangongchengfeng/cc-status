Status: done

# 实现失败恢复与部分提交语义

## What to build

补齐 `sync` 的失败恢复能力，使 client 在网络波动、服务端异常和文件变更边界下仍保持正确的状态语义。该切片需要实现指数退避重试、整批失败判定、按成功部分提交、严格重扫、文件截断后从头重扫，以及“部分失败即非 0 退出”的命令行为。完成后，client 对成功批次和失败批次的处理边界应清晰且可验证。

## Acceptance criteria

- [x] 网络错误、超时和 HTTP 5xx 会触发指数退避重试，最多 3 次；若仍失败则该批不写入 `reported_ids`
- [x] 当服务端返回 HTTP 200 但业务码异常，或 `accepted_count + duplicate_count` 小于本批记录数时，该批按失败处理且不推进任何状态
- [x] 已明确成功的批次会被保留为已提交结果，失败批次与关联文件在下次运行时会按严格重扫语义重新处理
- [x] 检测到文件截断、替换或行号回退时，client 会把该文件视为重置并从头重扫；本次运行若存在部分失败则命令以非 0 退出

## Blocked by

- [03-sync-happy-path-with-state-advance.md](./03-sync-happy-path-with-state-advance.md)
