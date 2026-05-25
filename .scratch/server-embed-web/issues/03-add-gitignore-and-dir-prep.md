Status: done

# 添加 Git 忽略与目录准备

## What to build

在 server/.gitignore 中添加 internal/ui/dist/ 忽略规则，创建必要的目录结构（如 internal/ui/ 下的 .gitkeep），确保构建产物不会被提交到 Git。

## Acceptance criteria

- [x] server/.gitignore 已添加 internal/ui/dist/ 规则。
- [x] server/internal/ui/ 目录存在（可通过 .gitkeep 占位）。
- [x] web/.gitignore 已有 dist/（已验证）。

## Blocked by

None - can start immediately
